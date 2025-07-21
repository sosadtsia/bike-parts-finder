package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sosadtsia/bike-parts-finder/pkg/database"
	"github.com/sosadtsia/bike-parts-finder/pkg/kafka"
	"github.com/sosadtsia/bike-parts-finder/pkg/models"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "CONSUMER: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting Bike Parts Finder Consumer")

	// Create context that listens for termination signals
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Println("Received shutdown signal")
		cancel()
	}()

	// Initialize database connection
	db, err := database.NewPostgresClient()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Kafka consumer for scrape results
	consumer, err := kafka.NewConsumer("scrape_results")
	if err != nil {
		logger.Fatalf("Failed to initialize Kafka consumer: %v", err)
	}

	// Process scrape results
	for {
		select {
		case <-ctx.Done():
			logger.Println("Shutting down gracefully...")
			consumer.Close()
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
			var result models.ScrapeResult
			if err := json.Unmarshal(msg.Value, &result); err != nil {
				logger.Printf("Error parsing scrape result: %v", err)
				continue
			}

			logger.Printf("Received %d parts from scrape of %s", len(result.Parts), result.URL)

			// Store parts in database
			for _, part := range result.Parts {
				if err := db.StorePart(ctx, part); err != nil {
					logger.Printf("Error storing part %s: %v", part.ID, err)
					continue
				}
			}

			logger.Printf("Successfully stored %d parts from %s", len(result.Parts), result.URL)

			// Commit the message offset
			if err := consumer.CommitMessage(msg); err != nil {
				logger.Printf("Error committing message offset: %v", err)
			}
		}
	}
}
