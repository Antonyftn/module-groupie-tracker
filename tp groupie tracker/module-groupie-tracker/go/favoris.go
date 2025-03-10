package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// NE PAS REDECLARER `UserFavorites` NI `FavoritesManager` car ils sont déjà dans `main.go`.

// Ajout des méthodes pour gérer les favoris
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
