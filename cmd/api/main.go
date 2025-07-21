package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sosadtsia/bike-parts-finder/pkg/api/handlers"
	"github.com/sosadtsia/bike-parts-finder/pkg/api/middleware"
	"github.com/sosadtsia/bike-parts-finder/pkg/cache"
	"github.com/sosadtsia/bike-parts-finder/pkg/database"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "API: ", log.LstdFlags|log.Lshortfile)
	logger.Println("Starting Bike Parts Finder API server")

	// Initialize database connection
	db, err := database.NewPostgresClient()
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Redis cache
	cacheClient, err := cache.NewRedisClient()
	if err != nil {
		logger.Printf("Warning: Failed to connect to Redis: %v", err)
		// Continue without cache
	} else {
		defer cacheClient.Close()
	}

	// Initialize router
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.Logging(logger))
	router.Use(middleware.CORS)

	// Initialize handlers
	partHandler := handlers.NewPartHandler(db, cacheClient)

	// Register routes
	router.HandleFunc("/api/parts", partHandler.GetAllParts).Methods("GET")
	router.HandleFunc("/api/parts/{id}", partHandler.GetPartByID).Methods("GET")
	router.HandleFunc("/api/parts/search", partHandler.SearchParts).Methods("GET")

	// Serve WebAssembly content
	fs := http.FileServer(http.Dir("./web/dist"))
	router.PathPrefix("/").Handler(fs)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Printf("Server listening on port %s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
