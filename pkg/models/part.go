package models

import "time"

// Part represents a bicycle part with its metadata
type Part struct {
	ID          string    `json:"id" db:"id"`
	Brand       string    `json:"brand" db:"brand"`
	Model       string    `json:"model" db:"model"`
	Year        int       `json:"year" db:"year"`
	Name        string    `json:"name" db:"name"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	SourceURL   string    `json:"source_url" db:"source_url"`
	Price       float64   `json:"price" db:"price"`
	Currency    string    `json:"currency" db:"currency"`
	InStock     bool      `json:"in_stock" db:"in_stock"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// SearchParams represents the parameters for searching parts
type SearchParams struct {
	Brand    string `json:"brand" form:"brand"`
	Model    string `json:"model" form:"model"`
	Year     int    `json:"year" form:"year"`
	Category string `json:"category" form:"category"`
	Page     int    `json:"page" form:"page" default:"1"`
	Limit    int    `json:"limit" form:"limit" default:"20"`
}

// ScrapeRequest represents a request to scrape parts data
type ScrapeRequest struct {
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
	Category string `json:"category,omitempty"`
	Source   string `json:"source,omitempty"` // Optional source to scrape from
}

// ScrapeResult represents the result of a scrape operation
type ScrapeResult struct {
	RequestID string  `json:"request_id"`
	Source    string  `json:"source"`
	Parts     []Part  `json:"parts"`
	Error     *string `json:"error,omitempty"`
	Count     int     `json:"count"`
	Timestamp int64   `json:"timestamp"`
}
