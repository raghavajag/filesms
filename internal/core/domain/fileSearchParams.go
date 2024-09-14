package domain

import (
	"time"
)

type FileSearchParams struct {
	Query    string
	FileType string
	FromDate time.Time
	ToDate   time.Time
	SortBy   string
	SortDir  string
	Limit    int
	Offset   int
}
