package main

import (
	"fmt"
	"math"
	"time"
)

type PaginationOptions struct {
	Page     int
	PageSize int
}

type PaginatedResult struct {
	Items       []Song
	TotalItems  int
	CurrentPage int
	PageSize    int
	TotalPages  int
	HasNext     bool
	HasPrev     bool
}

func NewCollection() *MusicCollection {
	return &MusicCollection{
		songs: make([]Song, 0),
	}
}

func (mc *MusicCollection) AddASong(song Song) {
	mc.songs = append(mc.songs, song)
}

func paginate(items []Song, opts PaginationOptions) PaginatedResult {
	totalItems := len(items)
	totalPages := int(math.Ceil(float64(totalItems) / float64(opts.PageSize)))

	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.Page > totalPages {
		opts.Page = totalPages
	}

	startIndex := (opts.Page - 1) * opts.PageSize
	endIndex := startIndex + opts.PageSize

	if endIndex > totalItems {
		endIndex = totalItems
	}

	var pageItems []Song
	if startIndex < totalItems {
		pageItems = items[startIndex:endIndex]
	} else {
		pageItems = []Song{}
	}

	return PaginatedResult{
		Items:       pageItems,
		TotalItems:  totalItems,
		CurrentPage: opts.Page,
		PageSize:    opts.PageSize,
		TotalPages:  totalPages,
		HasNext:     opts.Page < totalPages,
		HasPrev:     opts.Page > 1,
	}
}

func (mc *MusicCollection) SearchWithPagination(query string, pagination PaginationOptions) PaginatedResult {
	searchResults := mc.Search(query)
	return paginate(searchResults, pagination)
}

func (mc *MusicCollection) FilterWithPagination(filter Filter, pagination PaginationOptions) PaginatedResult {
	filterResults := mc.Filter(filter)
	return paginate(filterResults, pagination)
}

func (mc *MusicCollection) SearchFilterWithPagination(query string, filter Filter, pagination PaginationOptions) PaginatedResult {
	results := mc.SearchAndFilter(query, filter)
	return paginate(results, pagination)
}

func mainPagination() {

	collection := NewMusicCollection()
	for i := 1; i <= 50; i++ {
		collection.AddSong(Song{
			Title:       fmt.Sprintf("Song %d", i),
			Artist:      fmt.Sprintf("Artist %d", (i%5)+1),
			Album:       fmt.Sprintf("Album %d", (i%10)+1),
			Duration:    180 + i*30,
			Genre:       fmt.Sprintf("Genre %d", (i%3)+1),
			ReleaseDate: time.Now().AddDate(0, -i, 0),
		})
	}

	paginationOpts := PaginationOptions{
		Page:     2,
		PageSize: 10,
	}

	searchResults := collection.SearchWithPagination("Song", paginationOpts)

	fmt.Printf("Résultats paginés (Page %d/%d, %d éléments par page)\n",
		searchResults.CurrentPage, searchResults.TotalPages, searchResults.PageSize)
	fmt.Printf("Total des éléments: %d\n", searchResults.TotalItems)
	fmt.Printf("Page précédente ? %v, Page suivante ? %v\n\n",
		searchResults.HasPrev, searchResults.HasNext)

	for i, song := range searchResults.Items {
		fmt.Printf("%d. %s par %s\n", (searchResults.CurrentPage-1)*searchResults.PageSize+i+1,
			song.Title, song.Artist)
	}

	filter := Filter{
		MinDuration: func(i int) *int { return &i }(300),
	}

	filteredPageResults := collection.SearchFilterWithPagination("Song", filter, paginationOpts)

	fmt.Printf("\nRésultats filtrés et paginés (Page %d/%d)\n",
		filteredPageResults.CurrentPage, filteredPageResults.TotalPages)
	for i, song := range filteredPageResults.Items {
		fmt.Printf("%d. %s par %s (Durée: %d sec)\n",
			(filteredPageResults.CurrentPage-1)*filteredPageResults.PageSize+i+1,
			song.Title, song.Artist, song.Duration)
	}
}
