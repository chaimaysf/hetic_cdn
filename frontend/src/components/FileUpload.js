import React, { useState } from "react";
import axios from "axios";

const FileUpload = ({ onUploadSuccess }) => {
  const [selectedFile, setSelectedFile] = useState(null);
  const [progress, setProgress] = useState(0);
  const [message, setMessage] = useState("");

  const handleFileChange = (event) => {
    setSelectedFile(event.target.files[0]);
  };

  const handleUpload = async () => {
    if (!selectedFile) {
      setMessage("❌ Sélectionnez un fichier avant d'envoyer.");
      return;
    }

    const formData = new FormData();
    formData.append("file", selectedFile);

    try {
      const response = await axios.post("https://localhost:8080/upload", formData, {
        headers: { "Content-Type": "multipart/form-data" },
        onUploadProgress: (progressEvent) => {
          const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          setProgress(percentCompleted);
        },
      });

      setMessage("✅ Fichier uploadé avec succès !");
      setSelectedFile(null);
      setProgress(0);
      onUploadSuccess(); // 🔄 Met à jour la liste des fichiers
    } catch (error) {
      console.error("❌ Erreur lors de l'upload :", error);
      setMessage("❌ Erreur lors de l'upload : " + (error.response?.data || error.message));
    }
  };

  return (
    <div>
      <h3>📤 Uploader un fichier</h3>
      <input type="file" onChange={handleFileChange} />
      <button onClick={handleUpload}>Upload</button>
      {progress > 0 && <progress value={progress} max="100">{progress}%</progress>}
      {message && <p>{message}</p>}
    </div>
  );
};

export default FileUpload;
