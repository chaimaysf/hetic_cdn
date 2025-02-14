package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/chaim/hetic-cdn/models"
	"github.com/chaim/hetic-cdn/db"
	// "github.com/go-redis/redis/v8" // ✅ Ajout du package Redis
)


// 📂 Récupérer tous les fichiers
func GetFiles(w http.ResponseWriter, r *http.Request) {
	// 🔥 Vérifie d'abord si les fichiers sont déjà stockés dans Redis
	cachedFiles, err := db.RedisClient.Get(db.Ctx, "files_list").Result()
	if err == nil {
		// ✅ Si Redis contient les fichiers, on les renvoie directement
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cachedFiles))
		return
	}

	// 🛑 Si Redis ne contient pas les fichiers, on interroge la base de données
	rows, err := db.DB.Query("SELECT id, user_id, name, path, is_folder, parent_id, created_at, size, folder_id FROM files")
	if err != nil {
		log.Println("❌ Erreur de requête SQL :", err)
		http.Error(w, "Erreur de requête SQL", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []models.File
	for rows.Next() {
		var file models.File
		err := rows.Scan(&file.ID, &file.UserID, &file.Name, &file.Path, &file.IsFolder, &file.ParentID, &file.CreatedAt, &file.Size, &file.FolderID)
		if err != nil {
			log.Println("❌ Erreur de lecture des fichiers :", err)
			http.Error(w, "Erreur de lecture des fichiers", http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}

	// ✅ Convertir les données en JSON et les stocker dans Redis pour 10 minutes
	jsonData, _ := json.Marshal(files)
	db.RedisClient.Set(db.Ctx, "files_list", jsonData, 600000000000) // 10 min

	// ✅ Retourne les fichiers récupérés
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}


// 📂 Fonction pour uploader un fichier avec validation
func UploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(50 << 20) // Augmentation de la limite à 50 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "❌ Erreur lors du chargement du fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 📂 Définition du chemin du fichier
	uploadDir := "uploads"
	os.MkdirAll(uploadDir, os.ModePerm)
	fileName := handler.Filename
	filePath := filepath.Join(uploadDir, fileName)

	// 🛑 Vérifier l’extension du fichier
	allowedExtensions := []string{".png", ".jpg", ".jpeg", ".gif", ".pdf", ".txt", ".zip"}
	ext := filepath.Ext(fileName)
	isAllowed := false
	for _, validExt := range allowedExtensions {
		if ext == validExt {
			isAllowed = true
			break
		}
	}
	if !isAllowed {
		http.Error(w, "❌ Type de fichier non autorisé", http.StatusBadRequest)
		return
	}

	// 📌 Vérifier si le fichier existe déjà sur le disque
	if _, err := os.Stat(filePath); err == nil {
		http.Error(w, "⚠️ Fichier déjà existant", http.StatusConflict)
		return
	}

	// 📌 Vérifier si le fichier est déjà en base de données
	var existingID int
	err = db.DB.QueryRow("SELECT id FROM files WHERE name = $1", fileName).Scan(&existingID)
	if err == nil {
		http.Error(w, "⚠️ Fichier déjà enregistré", http.StatusConflict)
		return
	}

	// 📤 Enregistrement du fichier sur le disque
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "❌ Erreur interne", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	io.Copy(dst, file)

	// 🗃️ Sauvegarde en base de données
	var fileID int
	err = db.DB.QueryRow("INSERT INTO files (name, path, size) VALUES ($1, $2, $3) RETURNING id",
		fileName, filePath, handler.Size).Scan(&fileID)
	if err != nil {
		http.Error(w, "❌ Erreur base de données", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "✅ Fichier uploadé avec succès",
		"id":      fileID,
		"path":    filePath,
	})
}

// 📥 Télécharger un fichier via son ID
func DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 🔍 Vérifier si le fichier existe en base
	var filePath, fileName string
	err = db.DB.QueryRow("SELECT path, name FROM files WHERE id = $1", id).Scan(&filePath, &fileName)
	if err != nil {
		http.Error(w, "Fichier non trouvé en base", http.StatusNotFound)
		return
	}

	// 🔍 Vérifier si le fichier existe physiquement sur le disque
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Fichier introuvable sur le disque", http.StatusNotFound)
		return
	}

	// 📥 Envoi du fichier
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
// ❌ Supprimer un fichier par son ID avec une clé secrète
func DeleteFile(w http.ResponseWriter, r *http.Request) {
	// 🔑 Vérifier la clé secrète envoyée par le client
	secretKey := r.Header.Get("X-Secret-Key")
	if secretKey != "monSuperMotDePasse" { // Remplace par un vrai mot de passe sécurisé
		http.Error(w, "Accès interdit", http.StatusUnauthorized)
		return
	}

	// 🔍 Récupérer l'ID depuis l'URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID du fichier manquant", http.StatusBadRequest)
		return
	}

	// 🔢 Convertir en entier
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("❌ ID invalide :", err)
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// 🗃️ Vérifier si le fichier est un dossier contenant d'autres fichiers
	var isFolder bool
	err = db.DB.QueryRow("SELECT is_folder FROM files WHERE id = $1", id).Scan(&isFolder)
	if err != nil {
		log.Println("❌ Fichier non trouvé :", err)
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	if isFolder {
		// 🗑️ Supprimer les fichiers enfants d'un dossier
		_, err = db.DB.Exec("DELETE FROM files WHERE parent_id = $1", id)
		if err != nil {
			log.Println("❌ Erreur suppression fichiers enfants :", err)
			http.Error(w, "Erreur suppression fichiers enfants", http.StatusInternalServerError)
			return
		}
	}

	// 📂 Récupérer le chemin du fichier
	var filePath string
	err = db.DB.QueryRow("SELECT path FROM files WHERE id = $1", id).Scan(&filePath)
	if err != nil {
		http.Error(w, "Fichier non trouvé", http.StatusNotFound)
		return
	}

	// 🗑️ Supprimer le fichier du disque
	if filePath != "" {
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			log.Println("❌ Erreur suppression fichier disque :", err)
			http.Error(w, "Erreur suppression fichier", http.StatusInternalServerError)
			return
		}
	}

	// 🗃️ Supprimer l'entrée de la base de données
	_, err = db.DB.Exec("DELETE FROM files WHERE id = $1", id)
	if err != nil {
		log.Println("❌ Erreur suppression en base :", err)
		http.Error(w, "Erreur suppression base de données", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Fichier supprimé avec succès :", filePath)
	w.Write([]byte("✅ Fichier supprimé avec succès"))
}
