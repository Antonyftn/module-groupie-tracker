MusicHub
📌 À propos
MusicHub est une application web de gestion de bibliothèque musicale permettant aux utilisateurs de découvrir, rechercher et sauvegarder leurs morceaux préférés. L'application offre une interface intuitive et des fonctionnalités avancées de recherche et de filtrage.
🎵 Fonctionnalités principales

Recherche avancée de musiques
Filtrage par genre, durée et date de sortie
Système de pagination
Gestion des favoris avec persistance
Interface responsive
Catégorisation des musiques

🛠 Technologies utilisées

Go (Backend)
React (Frontend)
Tailwind CSS (Styling)
JSON (Persistance des données)

📋 Prérequis

Go 1.18 ou supérieur
Node.js 16.x ou supérieur
npm ou yarn

🚀 Installation
Backend
bashCopy# Cloner le repository
git clone https://github.com/votre-username/musichub.git
cd musichub

🌐 Structure des routes
Frontend
Copy/                   # Page d'accueil
/collection         # Liste des musiques
/collection/:id     # Détails d'une musique
/categories         # Liste des catégories
/categories/:id     # Musiques par catégorie
/favoris           # Gestion des favoris
/recherche         # Recherche avancée
/about             # À propos
Backend (API)
CopyGET    /api/songs              # Liste des musiques (avec pagination)
GET    /api/songs/:id          # Détails d'une musique
GET    /api/songs/search       # Recherche de musiques
GET    /api/categories         # Liste des catégories
GET    /api/categories/:id     # Musiques par catégorie
POST   /api/favorites          # Ajouter aux favoris
GET    /api/favorites          # Liste des favoris
DELETE /api/favorites/:id      # Supprimer des favoris
📡 Endpoints API utilisés
Recherche et filtrage
CopyGET /api/songs/search
Paramètres:
- query (string): Terme de recherche
- page (int): Numéro de page
- pageSize (int): Nombre d'éléments par page
- genre (string, optionnel): Filtre par genre
- minDuration (int, optionnel): Durée minimum en secondes
- maxDuration (int, optionnel): Durée maximum en secondes
- fromDate (string, optionnel): Date de sortie minimum (YYYY-MM-DD)
- toDate (string, optionnel): Date de sortie maximum (YYYY-MM-DD)
Gestion des favoris
CopyPOST /api/favorites
Body:
{
    "songId": "string",
    "userId": "string"
}

DELETE /api/favorites/:id
Paramètres:
- id: ID de la chanson à supprimer des favoris
📁 Structure des fichiers
Copymusichub/
├── backend/
│   ├── main.go
│   ├── api/
│   ├── models/
│   └── data/
├── frontend/
│   ├── src/
│   │   ├── components/
│   │   ├── pages/
│   │   ├── services/
│   │   └── utils/
│   └── public/
└── README.md