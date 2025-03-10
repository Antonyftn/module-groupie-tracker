package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/zmb3/spotify"
)

// ---- STRUCTURES ----

// Structure pour une chanson
type Song struct {
	Title       string
	Artist      string
	Album       string
	Duration    int
	Genre       string
	ReleaseDate time.Time
}

// Structure pour stocker les favoris d'un utilisateur
type UserFavorites struct {
	UserID    string    `json:"userId"`
	Songs     []Song    `json:"songs"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Structure de gestion des favoris
type FavoritesManager struct {
	baseDir string
}

// Structure des options de pagination
type PaginationOptions struct {
	Page     int
	PageSize int
}

// Structure des résultats paginés
type PaginatedResult struct {
	Items       []Song
	TotalItems  int
	CurrentPage int
	PageSize    int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
}

// ---- FONCTIONS DE GESTION DES FAVORIS ----

func NewFavoritesManager(baseDir string) (*FavoritesManager, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("erreur création répertoire favoris: %w", err)
	}
	return &FavoritesManager{baseDir: baseDir}, nil
}

func (fm *FavoritesManager) getFavoritesPath(userID string) string {
	return filepath.Join(fm.baseDir, fmt.Sprintf("favorites_%s.json", userID))
}

// Removed duplicate LoadFavorites method

// ---- GESTION DE L'AUTHENTIFICATION SPOTIFY ----

var (
	auth  = spotify.NewAuthenticator("http://localhost:8080/callback", spotify.ScopePlaylistReadPrivate)
	ch    = make(chan *spotify.Client)
	state = "random-string-for-security"
)

func StartAuth() {
	clientID := os.Getenv("d94076b3d4174dfca04516c874590d20")
	clientSecret := os.Getenv("8cb5a32df82541d59438409bbcb79bb5")

	if clientID == "" || clientSecret == "" {
		log.Fatal("SPOTIFY_CLIENT_ID ou SPOTIFY_CLIENT_SECRET non défini")
	}

	auth.SetAuthInfo(clientID, clientSecret)

	http.HandleFunc("/callback", completeAuth)
	go http.ListenAndServe(":8080", nil)

	url := auth.AuthURL(state)
	fmt.Printf("Veuillez ouvrir ce lien pour vous authentifier :\n%s\n", url)

	client := <-ch

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatalf("Impossible de récupérer l'utilisateur: %v", err)
	}
	fmt.Printf("Connecté en tant que : %s (%s)\n", user.DisplayName, user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Impossible d'obtenir le token", http.StatusForbidden)
		log.Fatalf("Erreur récupération token: %v", err)
		return
	}
	if r.FormValue("state") != state {
		http.NotFound(w, r)
		log.Fatal("État invalide")
		return
	}
	client := auth.NewClient(tok)
	fmt.Fprintln(w, "Authentification réussie! Vous pouvez revenir à l'application.")
	ch <- &client
}

// ---- HANDLERS HTTP ----

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Groupie Tracker !")
}

func collectionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Page Collection")
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Page Recherche")
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Page Filtrage")
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("templates/" + tmpl + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

// ---- POINT D'ENTRÉE ----

func main() {
	// Gestion des favoris
	_, err := NewFavoritesManager("./data/favorites")
	if err != nil {
		fmt.Printf("Erreur initialisation gestionnaire favoris: %v\n", err)
		return
	}
	fmt.Println("Gestion des favoris initialisée")

	// Démarrer l'authentification Spotify
	StartAuth()

	// Routes HTTP
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/collection", collectionHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/filter", filterHandler)

	fmt.Println("Serveur Web démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
