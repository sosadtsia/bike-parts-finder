package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sosadtsia/bike-parts-finder/pkg/cache"
	"github.com/sosadtsia/bike-parts-finder/pkg/database"
)

// PartHandler handles part-related API requests
type PartHandler struct {
	db    *database.PostgresClient
	cache *cache.RedisClient
}

// NewPartHandler creates a new part handler
func NewPartHandler(db *database.PostgresClient, cache *cache.RedisClient) *PartHandler {
	return &PartHandler{
		db:    db,
		cache: cache,
	}
}

// GetAllParts returns all parts
func (h *PartHandler) GetAllParts(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	limit := 20
	offset := 0

	// Get parts from database
	parts, err := h.db.GetParts(r.Context(), offset, limit)
	if err != nil {
		http.Error(w, "Error fetching parts", http.StatusInternalServerError)
		return
	}

	// Return parts as JSON
	if err := json.NewEncoder(w).Encode(parts); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// GetPartByID returns a part by ID
func (h *PartHandler) GetPartByID(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get path parameters
	vars := mux.Vars(r)
	id := vars["id"]

	// Try to get part from cache first
	part, err := h.cache.GetCachedPart(r.Context(), id)
	if err == nil {
		// Return cached part as JSON
		if err := json.NewEncoder(w).Encode(part); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
		return
	}

	// Get part from database
	part, err = h.db.GetPartByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Part not found", http.StatusNotFound)
		return
	}

	// Cache the part for future requests
	if h.cache != nil {
		go h.cache.CachePart(r.Context(), part)
	}

	// Return part as JSON
	if err := json.NewEncoder(w).Encode(part); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// SearchParts searches for parts
func (h *PartHandler) SearchParts(w http.ResponseWriter, r *http.Request) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Get query parameters
	query := r.URL.Query()
	q := query.Get("q")
	brand := query.Get("brand")
	category := query.Get("category")
	limit := 20
	offset := 0

	// Search for parts in database
	parts, err := h.db.SearchParts(r.Context(), q, brand, category, offset, limit)
	if err != nil {
		http.Error(w, "Error searching parts", http.StatusInternalServerError)
		return
	}

	// Return parts as JSON
	if err := json.NewEncoder(w).Encode(parts); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
