package models

import (
	"fmt"
	"math"
)

type PageInfo struct {
	TotalResults int         `json:"totalResults"`
	MaxResults   int         `json:"maxResults"`
	Paginations  Paginations `json:"-"`
}

type Pagination struct {
	Start        int
	CurrentStart int
	Page         int
	URL          string
}

type Paginations []*Pagination

func (page *PageInfo) ConfigPaginations(pageURL string, currentStart int) {
	paginations := Paginations{}
	pageTotal := int(math.Ceil(float64(page.TotalResults) / float64(page.MaxResults)))
	for i := 0; i < pageTotal; i++ {
		start := i * page.MaxResults
		page := &Pagination{
			Start:        start,
			CurrentStart: currentStart,
			Page:         i + 1,
			URL:          fmt.Sprintf("%sstart=%d&maxResults=%d", pageURL, start, page.MaxResults),
		}
		paginations = append(paginations, page)
	}
	page.Paginations = paginations
}
