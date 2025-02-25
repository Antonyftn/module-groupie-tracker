package main

import (
	"fmt"
	"strings"
)

type MusicCollection struct {
	songs []Song
}

func (mc *MusicCollection) SearchByField(query string, field string) []Song {
	query = strings.ToLower(query)
	results := make([]Song, 0)

	for _, song := range mc.songs {
		switch strings.ToLower(field) {
		case "title":
			if strings.Contains(strings.ToLower(song.Title), query) {
				results = append(results, song)
			}
		case "artist":
			if strings.Contains(strings.ToLower(song.Artist), query) {
				results = append(results, song)
			}
		case "album":
			if strings.Contains(strings.ToLower(song.Album), query) {
				results = append(results, song)
			}
		}
	}

	return results
}

func run() {
	collection := NewMusicCollection()

	collection.AddSong(Song{
		Title:    "Bohemian Rhapsody",
		Artist:   "Queen",
		Album:    "A Night at the Opera",
		Duration: 354,
	})
	collection.AddSong(Song{
		Title:    "Hotel California",
		Artist:   "Eagles",
		Album:    "Hotel California",
		Duration: 391,
	})

	query := "queen"
	results := collection.Search(query)
	fmt.Printf("Résultats pour '%s':\n", query)
	for _, song := range results {
		fmt.Printf("- %s par %s (Album: %s)\n", song.Title, song.Artist, song.Album)
	}

	fieldQuery := "hotel"
	fieldResults := collection.SearchByField(fieldQuery, "album")
	fmt.Printf("\nRésultats pour '%s' dans les albums:\n", fieldQuery)
	for _, song := range fieldResults {
		fmt.Printf("- %s par %s (Album: %s)\n", song.Title, song.Artist, song.Album)
	}
}
