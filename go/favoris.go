package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type UserFavorites struct {
	UserID    string    `json:"userId"`
	Songs     []Song    `json:"songs"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FavoritesManager struct {
	baseDir string
}

func NewFavoritesManager(baseDir string) (*FavoritesManager, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("erreur création répertoire favoris: %w", err)
	}
	return &FavoritesManager{baseDir: baseDir}, nil
}

func (fm *FavoritesManager) getFavoritesPath(userID string) string {
	return filepath.Join(fm.baseDir, fmt.Sprintf("favorites_%s.json", userID))
}

func (fm *FavoritesManager) LoadFavorites(userID string) (*UserFavorites, error) {
	path := fm.getFavoritesPath(userID)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &UserFavorites{
			UserID:    userID,
			Songs:     make([]Song, 0),
			UpdatedAt: time.Now(),
		}, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture favoris: %w", err)
	}

	var favorites UserFavorites
	if err := json.Unmarshal(data, &favorites); err != nil {
		return nil, fmt.Errorf("erreur décodage favoris: %w", err)
	}

	return &favorites, nil
}

func (fm *FavoritesManager) SaveFavorites(favorites *UserFavorites) error {
	favorites.UpdatedAt = time.Now()

	// Encoder en JSON
	data, err := json.MarshalIndent(favorites, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur encodage favoris: %w", err)
	}

	path := fm.getFavoritesPath(favorites.UserID)
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("erreur sauvegarde favoris: %w", err)
	}

	return nil
}

func (fm *FavoritesManager) AddToFavorites(userID string, song Song) error {
	favorites, err := fm.LoadFavorites(userID)
	if err != nil {
		return err
	}

	for _, s := range favorites.Songs {
		if s.Title == song.Title && s.Artist == song.Artist {
			return fmt.Errorf("cette chanson est déjà dans vos favoris")
		}
	}

	favorites.Songs = append(favorites.Songs, song)
	return fm.SaveFavorites(favorites)
}

func (fm *FavoritesManager) RemoveFromFavorites(userID string, songTitle, songArtist string) error {
	favorites, err := fm.LoadFavorites(userID)
	if err != nil {
		return err
	}

	newSongs := make([]Song, 0)
	found := false
	for _, s := range favorites.Songs {
		if s.Title != songTitle || s.Artist != songArtist {
			newSongs = append(newSongs, s)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("chanson non trouvée dans les favoris")
	}

	favorites.Songs = newSongs
	return fm.SaveFavorites(favorites)
}

func (fm *FavoritesManager) GetFavorites(userID string, pagination PaginationOptions) (PaginatedResult, error) {
	favorites, err := fm.LoadFavorites(userID)
	if err != nil {
		return PaginatedResult{}, err
	}

	return paginate(favorites.Songs, pagination), nil
}

func favoris() {
	favManager, err := NewFavoritesManager("./data/favorites")
	if err != nil {
		fmt.Printf("Erreur initialisation gestionnaire favoris: %v\n", err)
		return
	}

	userID := "user123"

	songs := []Song{
		{
			Title:       "Bohemian Rhapsody",
			Artist:      "Queen",
			Album:       "A Night at the Opera",
			Duration:    354,
			Genre:       "Rock",
			ReleaseDate: time.Date(1975, 11, 21, 0, 0, 0, 0, time.UTC),
		},
		{
			Title:       "Hotel California",
			Artist:      "Eagles",
			Album:       "Hotel California",
			Duration:    391,
			Genre:       "Rock",
			ReleaseDate: time.Date(1977, 2, 22, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, song := range songs {
		if err := favManager.AddToFavorites(userID, song); err != nil {
			fmt.Printf("Erreur ajout favori: %v\n", err)
			continue
		}
		fmt.Printf("Ajouté aux favoris: %s par %s\n", song.Title, song.Artist)
	}

	pagination := PaginationOptions{Page: 1, PageSize: 10}
	favorites, err := favManager.GetFavorites(userID, pagination)
	if err != nil {
		fmt.Printf("Erreur récupération favoris: %v\n", err)
		return
	}

	fmt.Printf("\nListe des favoris (%d total):\n", favorites.TotalItems)
	for i, song := range favorites.Items {
		fmt.Printf("%d. %s par %s (%s)\n", i+1, song.Title, song.Artist, song.Album)
	}

	if err := favManager.RemoveFromFavorites(userID, "Bohemian Rhapsody", "Queen"); err != nil {
		fmt.Printf("Erreur suppression favori: %v\n", err)
		return
	}
	fmt.Println("\nBohemian Rhapsody supprimé des favoris")

	favorites, err = favManager.GetFavorites(userID, pagination)
	if err != nil {
		fmt.Printf("Erreur récupération favoris: %v\n", err)
		return
	}

	fmt.Printf("\nListe des favoris mise à jour (%d total):\n", favorites.TotalItems)
	for i, song := range favorites.Items {
		fmt.Printf("%d. %s par %s (%s)\n", i+1, song.Title, song.Artist, song.Album)
	}
}
