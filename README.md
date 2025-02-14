# 🚀 Projet HETIC CDN

## 📌 Description du projet
Ce projet est un **CDN (Content Delivery Network)** développé en **Go**, utilisant **PostgreSQL**, **Redis**, **Prometheus** et **Grafana** pour le monitoring. Il inclut des fonctionnalités comme l’upload, le téléchargement, la suppression de fichiers et des mesures de sécurité comme un **firewall WAF**, un **système de bannissement IP**, et un **rate limiting**.

## 👥 **Membres du groupe**
- **Chaïma yousfi** 
- **Allaoua Lydia**
- **Wilton Ethan**

---

## 📦 **Technologies utilisées**
- **Backend**
  - Go (**Gorilla Mux** pour les routes)
  - PostgreSQL (**Base de données principale**)
  - Redis (**Cache, mais encore en développement**)
  - Prometheus (**Monitoring des métriques**)
  - Grafana (**Visualisation des métriques**)
  - **TLS/HTTPS** pour la sécurisation des échanges

- **Frontend**
  - React.js
  - Axios (pour les requêtes API)

---

## 📌 **Fonctionnalités implémentées ✅**
### **Backend**
✅ **Gestion des fichiers :**
  - 📤 Upload sécurisé de fichiers (taille max : 50MB, extensions limitées)
  - 📥 Téléchargement de fichiers via ID
  - ❌ Suppression sécurisée avec clé API (`X-Secret-Key`)

✅ **Sécurisation :**
  - 🔥 **Filtrage IP** (ban automatique après dépassement des limites)
  - 🚨 **Web Application Firewall (WAF)** (détection XSS & SQL Injection)
  - ⏳ **Rate Limiting** (évite les attaques DDoS)

✅ **Monitoring et métriques :**
  - 📊 **Prometheus** collecte les métriques d’usage du CDN
  - 📈 **Grafana** affiche les métriques en temps réel
  - 📍 **/metrics** retourne les statistiques d’usage du CDN

✅ **Base de données :**
  - PostgreSQL stocke **les fichiers** et **les dossiers**
  - Initialisation automatique des tables

### **Frontend**
✅ **Gestion des fichiers côté client :**
  - 🎨 Interface en **React.js**
  - 📤 Upload et 📥 téléchargement de fichiers
  - ❌ Suppression sécurisée (via clé API)

✅ **Requêtes API :**
  - Utilisation d’**Axios** pour interagir avec le backend



