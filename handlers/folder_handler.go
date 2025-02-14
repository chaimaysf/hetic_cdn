package handlers

import (
	"fmt"
	// "log"
	"net/http"

	"github.com/chaim/hetic-cdn/db"
)

func CreateFolder(userID int, name string, parentFolderID *int) error {
	query := `INSERT INTO folders (user_id, name, parent_folder) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, userID, name, parentFolderID)
	if err != nil {
		return fmt.Errorf("❌ Erreur lors de la création du dossier : %v", err)
	}
	fmt.Println("✅ Dossier créé avec succès !")
	return nil
}

func ListFolders(userID int) ([]string, error) {
	query := `SELECT name FROM folders WHERE user_id = $1`
	rows, err := db.DB.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("❌ Erreur lors de la récupération des dossiers : %v", err)
	}
	defer rows.Close()

	var folders []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		folders = append(folders, name)
	}
	return folders, nil
}

// Route HTTP pour lister les dossiers
func ListFoldersHandler(w http.ResponseWriter, r *http.Request) {
	userID := 1 // Remplace par l'ID du vrai utilisateur connecté
	folders, err := ListFolders(userID)
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des dossiers", http.StatusInternalServerError)
		return
	}

	for _, folder := range folders {
		fmt.Fprintf(w, "📁 %s\n", folder)
	}
}
