package main

import (
	"fmt"
	"strings"
	"time"
)

type Song struct {
	Title       string
	Artist      string
	Album       string
	Duration    int
	Genre       string
	ReleaseDate time.Time
}

type Filter struct {
	Genre       *string
	MinDuration *int
	MaxDuration *int
	FromDate    *time.Time
	ToDate      *time.Time
}

func NewMusicCollection() *MusicCollection {
	return &MusicCollection{
		songs: make([]Song, 0),
	}
}

func (mc *MusicCollection) AddSong(song Song) {
	mc.songs = append(mc.songs, song)
}

func (mc *MusicCollection) Search(query string) []Song {
	query = strings.ToLower(query)
	results := make([]Song, 0)

	for _, song := range mc.songs {
		if strings.Contains(strings.ToLower(song.Title), query) ||
			strings.Contains(strings.ToLower(song.Artist), query) ||
			strings.Contains(strings.ToLower(song.Album), query) {
			results = append(results, song)
		}
	}

	return results
}

func (mc *MusicCollection) Filter(filter Filter) []Song {
	results := make([]Song, 0)

	for _, song := range mc.songs {
		include := true

		if filter.Genre != nil && *filter.Genre != "" {
			if !strings.EqualFold(song.Genre, *filter.Genre) {
				include = false
			}
		}

		if filter.MinDuration != nil {
			if song.Duration < *filter.MinDuration {
				include = false
			}
		}

		if filter.MaxDuration != nil {
			if song.Duration > *filter.MaxDuration {
				include = false
			}
		}

		if filter.FromDate != nil {
			if song.ReleaseDate.Before(*filter.FromDate) {
				include = false
			}
		}

		if filter.ToDate != nil {
			if song.ReleaseDate.After(*filter.ToDate) {
				include = false
			}
		}

		if include {
			results = append(results, song)
		}
	}

	return results
}

func (mc *MusicCollection) SearchAndFilter(query string, filter Filter) []Song {
	searchResults := mc.Search(query)

	tempCollection := &MusicCollection{songs: searchResults}

	return tempCollection.Filter(filter)
}

func collection() {
	collection := NewMusicCollection()

	str := func(s string) *string { return &s }
	integer := func(i int) *int { return &i }
	date := func(t time.Time) *time.Time { return &t }

	collection.AddSong(Song{
		Title:       "Bohemian Rhapsody",
		Artist:      "Queen",
		Album:       "A Night at the Opera",
		Duration:    354,
		Genre:       "Rock",
		ReleaseDate: time.Date(1975, 11, 21, 0, 0, 0, 0, time.UTC),
	})
	collection.AddSong(Song{
		Title:       "Hotel California",
		Artist:      "Eagles",
		Album:       "Hotel California",
		Duration:    391,
		Genre:       "Rock",
		ReleaseDate: time.Date(1977, 2, 22, 0, 0, 0, 0, time.UTC),
	})

	filter := Filter{
		Genre:       str("Rock"),
		MinDuration: integer(360),
		FromDate:    date(time.Date(1976, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	results := collection.Filter(filter)
	fmt.Println("Résultats du filtrage :")
	for _, song := range results {
		fmt.Printf("- %s par %s (Genre: %s, Durée: %d sec, Date: %s)\n",
			song.Title, song.Artist, song.Genre, song.Duration,
			song.ReleaseDate.Format("2006-01-02"))
	}

	query := "hotel"
	combinedResults := collection.SearchAndFilter(query, filter)
	fmt.Printf("\nRésultats pour '%s' avec filtres :\n", query)
	for _, song := range combinedResults {
		fmt.Printf("- %s par %s (Genre: %s, Durée: %d sec, Date: %s)\n",
			song.Title, song.Artist, song.Genre, song.Duration,
			song.ReleaseDate.Format("2006-01-02"))
	}
}
