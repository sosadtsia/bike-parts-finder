package models

import "time"

// Part represents a bicycle part
type Part struct {
	ID          string    `json:"id"`
	Brand       string    `json:"brand"`
	Model       string    `json:"model"`
	Category    string    `json:"category"`
	SubCategory string    `json:"sub_category"`
	Price       float64   `json:"price"`
	MSRP        float64   `json:"msrp,omitempty"`
	Discount    float64   `json:"discount,omitempty"`
	Currency    string    `json:"currency"`
	InStock     bool      `json:"in_stock"`
	Rating      float64   `json:"rating,omitempty"`
	NumReviews  int       `json:"num_reviews,omitempty"`
	Description string    `json:"description"`
	Images      []string  `json:"images,omitempty"`
	URL         string    `json:"url"`
	Source      string    `json:"source,omitempty"`
	Specs       []Spec    `json:"specs,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// Spec represents a specification of a part
type Spec struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ScrapeRequest represents a request to scrape a URL for bike parts
type ScrapeRequest struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Source    string    `json:"source"`
	Timestamp time.Time `json:"timestamp"`
}

// ScrapeResult represents the result of a scraping operation
type ScrapeResult struct {
	RequestID string    `json:"request_id"`
	URL       string    `json:"url"`
	Parts     []Part    `json:"parts"`
	Timestamp time.Time `json:"timestamp"`
}
