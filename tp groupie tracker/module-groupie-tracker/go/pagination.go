package main

import (
	"math"
)

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
