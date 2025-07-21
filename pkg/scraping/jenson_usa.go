package scraping

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/sosadtsia/bike-parts-finder/pkg/models"
)

// JensonUSAScraper is a scraper for JensonUSA website
type JensonUSAScraper struct {
	baseURL string
}

// NewJensonUSAScraper creates a new JensonUSA scraper
func NewJensonUSAScraper() *JensonUSAScraper {
	return &JensonUSAScraper{
		baseURL: "https://www.jensonusa.com",
	}
}

// CanHandle checks if this scraper can handle the given URL
func (s *JensonUSAScraper) CanHandle(url string) bool {
	return strings.Contains(url, "jensonusa.com")
}

// Scrape scrapes bike parts from JensonUSA
func (s *JensonUSAScraper) Scrape(url string) ([]models.Part, error) {
	var parts []models.Part

	// Initialize the collector
	c := colly.NewCollector(
		colly.AllowedDomains("www.jensonusa.com", "jensonusa.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Category page scraper
	if strings.Contains(url, "/categories/") || strings.Contains(url, "/products/") {
		// For category pages, visit each product page
		c.OnHTML("div.product-tile", func(e *colly.HTMLElement) {
			productURL := e.ChildAttr("a.product-tile__image-link", "href")
			if productURL != "" {
				if !strings.HasPrefix(productURL, "http") {
					productURL = s.baseURL + productURL
				}
				// Visit the product page
				c.Visit(productURL)
			}
		})

		// Handle pagination
		c.OnHTML("a.pagination__link", func(e *colly.HTMLElement) {
			nextURL := e.Attr("href")
			if nextURL != "" && !strings.Contains(nextURL, "page=1") { // Avoid infinite loop
				c.Visit(s.baseURL + nextURL)
			}
		})
	}

	// Product page scraper
	c.OnHTML("div.product-details", func(e *colly.HTMLElement) {
		var part models.Part

		// Generate a unique ID for the part
		part.ID = uuid.New().String()

		// Get the product URL
		part.URL = e.Request.URL.String()
		part.Source = "JensonUSA"

		// Get the product name, which usually contains brand and model
		productName := e.ChildText("h1.product-details__name")
		nameParts := strings.SplitN(productName, " ", 2)
		if len(nameParts) > 0 {
			part.Brand = nameParts[0]
		}
		if len(nameParts) > 1 {
			part.Model = nameParts[1]
		}

		// Get the category from breadcrumbs
		e.ForEach("ol.breadcrumb li", func(_ int, el *colly.HTMLElement) {
			category := strings.TrimSpace(el.Text)
			if category != "Home" && category != "" {
				if part.Category == "" {
					part.Category = category
				} else {
					part.SubCategory = category
				}
			}
		})

		// Get the price
		priceText := e.ChildText("span.product-details__price--sale")
		if priceText == "" {
			priceText = e.ChildText("span.product-details__price")
		}

		// Clean up price and convert to float
		priceText = strings.TrimSpace(strings.ReplaceAll(priceText, "$", ""))
		price, err := strconv.ParseFloat(priceText, 64)
		if err == nil {
			part.Price = price
			part.Currency = "USD"
		}

		// Get MSRP if available
		msrpText := e.ChildText("span.product-details__price--msrp")
		msrpText = strings.TrimSpace(strings.ReplaceAll(msrpText, "$", ""))
		msrp, err := strconv.ParseFloat(msrpText, 64)
		if err == nil {
			part.MSRP = msrp

			// Calculate discount
			if part.MSRP > 0 {
				part.Discount = ((part.MSRP - part.Price) / part.MSRP) * 100
			}
		}

		// Check if in stock
		stockText := e.ChildText("div.product-details__stock")
		part.InStock = !strings.Contains(strings.ToLower(stockText), "out of stock")

		// Get the description
		part.Description = e.ChildText("div.product-details__description")

		// Get product images
		e.ForEach("div.product-details__image img", func(_ int, el *colly.HTMLElement) {
			imgURL := el.Attr("src")
			if imgURL != "" {
				if !strings.HasPrefix(imgURL, "http") {
					imgURL = "https:" + imgURL
				}
				part.Images = append(part.Images, imgURL)
			}
		})

		// Get specifications
		e.ForEach("table.specifications__table tr", func(_ int, el *colly.HTMLElement) {
			name := strings.TrimSpace(el.ChildText("td:first-child"))
			value := strings.TrimSpace(el.ChildText("td:last-child"))

			if name != "" && value != "" {
				part.Specs = append(part.Specs, models.Spec{
					Name:  name,
					Value: value,
				})
			}
		})

		// Set timestamps
		now := time.Now()
		part.CreatedAt = now
		part.UpdatedAt = now

		// Add the part to the results
		parts = append(parts, part)
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Request URL: %s failed with error: %s\n", r.Request.URL, err)
	})

	// Start the scraping
	err := c.Visit(url)
	if err != nil {
		return nil, fmt.Errorf("error starting the scraper: %w", err)
	}

	// Wait for scraping to finish
	c.Wait()

	return parts, nil
}
