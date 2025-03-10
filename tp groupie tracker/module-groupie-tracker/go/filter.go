package main

import (
	"strings"
	"time"
)

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
