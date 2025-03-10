package main

import (
	"fmt"
	"net/http"
)

type Resource struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Genre string `json:"genre"`
	Year  int    `json:"year"`
}

var resources = []Resource{
	{1, "Music 1", "Rock", 2020},
	{2, "Music 2", "Pop", 2019},
	{3, "Music 3", "Jazz", 2021},
}

func CollectionHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "collection", resources)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	var results []Resource
	for _, res := range resources {
		if query != "" && (res.Name == query || res.Genre == query) {
			results = append(results, res)
		}
	}
	renderTemplate(w, "search", results)
}

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	year := r.URL.Query().Get("year")
	var results []Resource
	for _, res := range resources {
		if (genre == "" || res.Genre == genre) && (year == "" || fmt.Sprintf("%d", res.Year) == year) {
			results = append(results, res)
		}
	}
	renderTemplate(w, "filter", results)
}
