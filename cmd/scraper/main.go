package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sosadtsia/bike-parts-finder/pkg/kafka"
	"github.com/sosadtsia/bike-parts-finder/pkg/models"
	"github.com/sosadtsia/bike-parts-finder/pkg/scraping"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "SCRAPER: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting Bike Parts Finder Scraper")

	// Create context that listens for termination signals
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Println("Received shutdown signal")
		cancel()
	}()

	// Initialize Kafka consumer for scrape requests
	consumer, err := kafka.NewConsumer("scrape_requests")
	if err != nil {
		logger.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	// Initialize Kafka producer for scrape results
	producer, err := kafka.NewProducer("scrape_results")
	if err != nil {
		logger.Fatalf("Failed to initialize Kafka producer: %v", err)
	}

	// Initialize scrapers
	jensonScraper := scraping.NewJensonUSAScraper()

	// Process scrape requests
	for {
		select {
		case <-ctx.Done():
			logger.Println("Shutting down gracefully...")
			consumer.Close()
			producer.Close()
			return
		default:
			// Poll for new messages
			msg, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// No message available
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Process message
			var request models.ScrapeRequest
			if err := json.Unmarshal(msg.Value, &request); err != nil {
				logger.Printf("Error parsing scrape request: %v", err)
				continue
			}

			logger.Printf("Received scrape request for: %s", request.URL)

			// Select the appropriate scraper
			var parts []models.Part

			if jensonScraper.CanHandle(request.URL) {
				parts, err = jensonScraper.Scrape(request.URL)
			} else {
				logger.Printf("No scraper available for URL: %s", request.URL)
				continue
			}

			if err != nil {
				logger.Printf("Error scraping %s: %v", request.URL, err)
				continue
			}

			// Send results to Kafka
			result := models.ScrapeResult{
				RequestID: request.ID,
				URL:       request.URL,
				Parts:     parts,
				Timestamp: time.Now(),
			}

			resultBytes, err := json.Marshal(result)
			if err != nil {
				logger.Printf("Error serializing scrape result: %v", err)
				continue
			}

			if err := producer.WriteMessage(resultBytes); err != nil {
				logger.Printf("Error sending scrape result to Kafka: %v", err)
			} else {
				logger.Printf("Successfully scraped %d parts from %s", len(parts), request.URL)
			}
		}
	}
}
