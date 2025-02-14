import React, { useState, useEffect } from "react";
import axios from "axios";
import FileUpload from "./components/FileUpload";
import FileList from "./components/FileList";

function App() {
  const [files, setFiles] = useState([]);

  // 🔄 Charger la liste des fichiers
  useEffect(() => {
    fetchFiles();
  }, []);

  const fetchFiles = async () => {
    try {
      const response = await axios.get("https://localhost:8080/files");
      setFiles(response.data);
    } catch (error) {
      console.error("❌ Erreur lors du chargement des fichiers :", error);
    }
  };

  // 📥 Télécharger un fichier
  const handleDownload = async (id, name) => {
    try {
      const response = await axios.get(`https://localhost:8080/download/${id}`, { responseType: "blob" });
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute("download", name);
      document.body.appendChild(link);
      link.click();
      link.remove();
    } catch (error) {
      console.error("❌ Erreur lors du téléchargement :", error);
    }
  };

  // ❌ Supprimer un fichier
  const handleDelete = async (id) => {
    try {
      await axios.delete(`https://localhost:8080/delete/${id}`);
      fetchFiles(); // Rafraîchir la liste
    } catch (error) {
      console.error("❌ Erreur lors de la suppression :", error);
    }
  };

  return (
    <div style={{ padding: "20px", fontFamily: "Arial, sans-serif" }}>
      <h2>📂 Gestionnaire de fichiers</h2>
      <FileUpload onUploadSuccess={fetchFiles} />
      <FileList files={files} onDownload={handleDownload} onDelete={handleDelete} />
    </div>
  );
}

export default App;
