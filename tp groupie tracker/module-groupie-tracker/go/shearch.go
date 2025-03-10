package main

import (
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
