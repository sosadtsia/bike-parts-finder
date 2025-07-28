package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	// Initialize router with strict slashes
	router := mux.NewRouter().StrictSlash(true)

	// Apply middleware
	router.Use(middleware.Logging(logger))
	router.Use(middleware.CORS)

	// Initialize handlers
	partHandler := handlers.NewPartHandler(db, cacheClient)

	// Health check endpoints
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	}).Methods("GET")

	router.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if err := db.Ping(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Database not ready: %v", err)
			return
		}

		// Check Redis if available
		if cacheClient != nil {
			if err := cacheClient.Ping(); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintf(w, "Cache not ready: %v", err)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ready")
	}).Methods("GET")

	// API v1 Routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/parts", partHandler.GetAllParts).Methods("GET")
	apiV1.HandleFunc("/parts/{id}", partHandler.GetPartByID).Methods("GET")
	apiV1.HandleFunc("/parts/search", partHandler.SearchParts).Methods("GET")

	// For backward compatibility
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

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	logger.Printf("Server listening on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server error: %v", err)
	}
}
