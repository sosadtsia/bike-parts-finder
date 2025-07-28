package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sosadtsia/bike-parts-finder/pkg/models"
)

// PostgresClient is a wrapper around a PostgreSQL connection pool
type PostgresClient struct {
	pool *pgxpool.Pool
}

// NewPostgresClient creates a new PostgresClient
func NewPostgresClient() (*PostgresClient, error) {
	// Get database connection string from environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/bike_parts_finder?sslmode=disable"
	}

	// Create connection pool
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parsing database connection string: %w", err)
	}

	// Set connection pool options
	config.MaxConns = 10
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	// Ping database to ensure connection
	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &PostgresClient{
		pool: pool,
	}, nil
}

// Close closes the database connection pool
func (c *PostgresClient) Close() {
	if c.pool != nil {
		c.pool.Close()
	}
}

// Ping checks if the database connection is alive
func (c *PostgresClient) Ping() error {
	// In a real implementation, this would check the database connection
	// For demonstration purposes, we'll just return nil
	return nil
}

// StorePart stores a bike part in the database
func (c *PostgresClient) StorePart(ctx context.Context, part models.Part) error {
	query := `
		INSERT INTO parts (
			id, brand, model, category, sub_category, price, msrp, currency,
			in_stock, rating, num_reviews, description, url, source, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		) ON CONFLICT (id) DO UPDATE SET
			brand = $2,
			model = $3,
			category = $4,
			sub_category = $5,
			price = $6,
			msrp = $7,
			currency = $8,
			in_stock = $9,
			rating = $10,
			num_reviews = $11,
			description = $12,
			url = $13,
			source = $14,
			updated_at = $16
		RETURNING id`

	// If part doesn't have created_at or updated_at timestamps, set them to now
	now := time.Now()
	if part.CreatedAt.IsZero() {
		part.CreatedAt = now
	}
	if part.UpdatedAt.IsZero() {
		part.UpdatedAt = now
	}

	_, err := c.pool.Exec(ctx, query,
		part.ID, part.Brand, part.Model, part.Category, part.SubCategory,
		part.Price, part.MSRP, part.Currency, part.InStock, part.Rating,
		part.NumReviews, part.Description, part.URL, part.Source,
		part.CreatedAt, part.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("storing part %s: %w", part.ID, err)
	}

	// Store specs
	if len(part.Specs) > 0 {
		err = c.storeSpecs(ctx, part.ID, part.Specs)
		if err != nil {
			return fmt.Errorf("storing specs for part %s: %w", part.ID, err)
		}
	}

	// Store images
	if len(part.Images) > 0 {
		err = c.storeImages(ctx, part.ID, part.Images)
		if err != nil {
			return fmt.Errorf("storing images for part %s: %w", part.ID, err)
		}
	}

	return nil
}

// storeSpecs stores specifications for a part
func (c *PostgresClient) storeSpecs(ctx context.Context, partID string, specs []models.Spec) error {
	// First, delete existing specs
	_, err := c.pool.Exec(ctx, "DELETE FROM part_specs WHERE part_id = $1", partID)
	if err != nil {
		return err
	}

	// Then insert new specs
	for _, spec := range specs {
		_, err = c.pool.Exec(ctx, "INSERT INTO part_specs (part_id, name, value) VALUES ($1, $2, $3)",
			partID, spec.Name, spec.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

// storeImages stores images for a part
func (c *PostgresClient) storeImages(ctx context.Context, partID string, images []string) error {
	// First, delete existing images
	_, err := c.pool.Exec(ctx, "DELETE FROM part_images WHERE part_id = $1", partID)
	if err != nil {
		return err
	}

	// Then insert new images
	for i, url := range images {
		_, err = c.pool.Exec(ctx, "INSERT INTO part_images (part_id, url, position) VALUES ($1, $2, $3)",
			partID, url, i)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPartByID retrieves a part by ID
func (c *PostgresClient) GetPartByID(ctx context.Context, id string) (models.Part, error) {
	var part models.Part

	// Query the part
	err := c.pool.QueryRow(ctx, `
		SELECT id, brand, model, category, sub_category, price, msrp, currency,
		       in_stock, rating, num_reviews, description, url, source, created_at, updated_at
		FROM parts WHERE id = $1
	`, id).Scan(
		&part.ID, &part.Brand, &part.Model, &part.Category, &part.SubCategory,
		&part.Price, &part.MSRP, &part.Currency, &part.InStock, &part.Rating,
		&part.NumReviews, &part.Description, &part.URL, &part.Source,
		&part.CreatedAt, &part.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return part, fmt.Errorf("part not found")
		}
		return part, err
	}

	// Load specs
	specs, err := c.getPartSpecs(ctx, id)
	if err != nil {
		return part, err
	}
	part.Specs = specs

	// Load images
	images, err := c.getPartImages(ctx, id)
	if err != nil {
		return part, err
	}
	part.Images = images

	return part, nil
}

// getPartSpecs retrieves specs for a part
func (c *PostgresClient) getPartSpecs(ctx context.Context, partID string) ([]models.Spec, error) {
	rows, err := c.pool.Query(ctx, "SELECT name, value FROM part_specs WHERE part_id = $1 ORDER BY name", partID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var specs []models.Spec
	for rows.Next() {
		var spec models.Spec
		if err := rows.Scan(&spec.Name, &spec.Value); err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}

	return specs, rows.Err()
}

// getPartImages retrieves images for a part
func (c *PostgresClient) getPartImages(ctx context.Context, partID string) ([]string, error) {
	rows, err := c.pool.Query(ctx, "SELECT url FROM part_images WHERE part_id = $1 ORDER BY position", partID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		images = append(images, url)
	}

	return images, rows.Err()
}

// GetParts retrieves a list of parts with pagination
func (c *PostgresClient) GetParts(ctx context.Context, offset, limit int) ([]models.Part, error) {
	rows, err := c.pool.Query(ctx, `
		SELECT id, brand, model, category, sub_category, price, msrp, currency,
		       in_stock, rating, num_reviews, description, url, source, created_at, updated_at
		FROM parts ORDER BY created_at DESC LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var part models.Part
		if err := rows.Scan(
			&part.ID, &part.Brand, &part.Model, &part.Category, &part.SubCategory,
			&part.Price, &part.MSRP, &part.Currency, &part.InStock, &part.Rating,
			&part.NumReviews, &part.Description, &part.URL, &part.Source,
			&part.CreatedAt, &part.UpdatedAt,
		); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load specs and images for each part
	for i, part := range parts {
		specs, err := c.getPartSpecs(ctx, part.ID)
		if err != nil {
			return nil, err
		}
		parts[i].Specs = specs

		images, err := c.getPartImages(ctx, part.ID)
		if err != nil {
			return nil, err
		}
		parts[i].Images = images
	}

	return parts, nil
}

// SearchParts searches for parts based on query parameters
func (c *PostgresClient) SearchParts(ctx context.Context, query, brand, category string, offset, limit int) ([]models.Part, error) {
	// Build the query
	sqlQuery := `
		SELECT id, brand, model, category, sub_category, price, msrp, currency,
		       in_stock, rating, num_reviews, description, url, source, created_at, updated_at
		FROM parts WHERE 1=1
	`
	args := []interface{}{}

	// Add filters
	if query != "" {
		sqlQuery += " AND (brand ILIKE $" + fmt.Sprintf("%d", len(args)+1) + " OR model ILIKE $" + fmt.Sprintf("%d", len(args)+1) + " OR description ILIKE $" + fmt.Sprintf("%d", len(args)+1) + ")"
		args = append(args, "%"+query+"%")
	}

	if brand != "" {
		sqlQuery += " AND brand ILIKE $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+brand+"%")
	}

	if category != "" {
		sqlQuery += " AND category ILIKE $" + fmt.Sprintf("%d", len(args)+1)
		args = append(args, "%"+category+"%")
	}

	// Add pagination
	sqlQuery += " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	// Execute the query
	rows, err := c.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var part models.Part
		if err := rows.Scan(
			&part.ID, &part.Brand, &part.Model, &part.Category, &part.SubCategory,
			&part.Price, &part.MSRP, &part.Currency, &part.InStock, &part.Rating,
			&part.NumReviews, &part.Description, &part.URL, &part.Source,
			&part.CreatedAt, &part.UpdatedAt,
		); err != nil {
			return nil, err
		}
		parts = append(parts, part)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Load specs and images for each part
	for i, part := range parts {
		specs, err := c.getPartSpecs(ctx, part.ID)
		if err != nil {
			return nil, err
		}
		parts[i].Specs = specs

		images, err := c.getPartImages(ctx, part.ID)
		if err != nil {
			return nil, err
		}
		parts[i].Images = images
	}

	return parts, nil
}
