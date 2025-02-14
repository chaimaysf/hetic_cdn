import React, { useState, useEffect } from "react";
import axios from "axios";

const FileList = () => {
  const [files, setFiles] = useState([]);
  const [message, setMessage] = useState("");

  // 🔄 Charger les fichiers au montage du composant
  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const response = await axios.get("https://localhost:8080/files");
      setFiles(response.data);
    } catch (error) {
      console.error("❌ Erreur lors de la récupération des fichiers", error);
    }
  };

  // 📥 Télécharger un fichier
  const handleDownload = async (id, name) => {
    try {
      const response = await axios.get(
        `https://localhost:8080/download/${id}`,
        { responseType: "blob" }
      );
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", name);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (error) {
      console.error("❌ Erreur lors du téléchargement", error);
    }
  };

  // ❌ Supprimer un fichier
  const handleDelete = async (id) => {
    try {
      await axios.delete(`https://localhost:8080/delete/${id}`, {
        headers: { "X-Secret-Key": "monSuperMotDePasse" }, // 🔑 Ajout de la clé secrète
      });
      setMessage("✅ Fichier supprimé avec succès !");
      fetchFiles(); // 🔄 Met à jour la liste après suppression
    } catch (error) {
      console.error("❌ Erreur lors de la suppression :", error);
      setMessage("❌ Erreur lors de la suppression.");
    }
  };

  return (
    <div>
      <h2>📄 Liste des fichiers</h2>
      {message && <p>{message}</p>}
      <ul>
        {files.map((file) => (
          <li key={file.id}>
            {file.name} - 📁 {file.path}
            <button onClick={() => handleDownload(file.id, file.name)}>
              📥 Télécharger
            </button>
            <button onClick={() => handleDelete(file.id)}>❌ Supprimer</button>
          </li>
        ))}
      </ul>
    </div>
  );
};
export default FileList;
