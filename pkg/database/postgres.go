package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/svosadtsia/bike-parts-finder/pkg/models"
)

// PostgresConfig holds the configuration for connecting to a PostgreSQL database
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

// Store represents a data store for bike parts
type Store struct {
	db *sql.DB
}

// NewPostgresStore creates a new store with a PostgreSQL connection
func NewPostgresStore(config PostgresConfig) (*Store, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Verify the connection is working
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &Store{db: db}, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	return s.db.Close()
}

// FindPartsByBrandModelYear finds bike parts matching the given criteria
func (s *Store) FindPartsByBrandModelYear(ctx context.Context, params models.SearchParams) ([]models.Part, error) {
	// Calculate offset based on page and limit
	offset := (params.Page - 1) * params.Limit

	query := `
		SELECT id, brand, model, year, name, category, description,
		image_url, source_url, price, currency, in_stock, created_at, updated_at
		FROM parts
		WHERE ($1 = '' OR brand ILIKE '%' || $1 || '%')
		AND ($2 = '' OR model ILIKE '%' || $2 || '%')
		AND ($3 = 0 OR year = $3)
		AND ($4 = '' OR category ILIKE '%' || $4 || '%')
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
	`

	rows, err := s.db.QueryContext(
		ctx,
		query,
		params.Brand,
		params.Model,
		params.Year,
		params.Category,
		params.Limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying parts: %w", err)
	}
	defer rows.Close()

	var parts []models.Part
	for rows.Next() {
		var part models.Part
		err := rows.Scan(
			&part.ID,
			&part.Brand,
			&part.Model,
			&part.Year,
			&part.Name,
			&part.Category,
			&part.Description,
			&part.ImageURL,
			&part.SourceURL,
			&part.Price,
			&part.Currency,
			&part.InStock,
			&part.CreatedAt,
			&part.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning part row: %w", err)
		}
		parts = append(parts, part)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating part rows: %w", err)
	}

	return parts, nil
}

// CountPartsByBrandModelYear counts bike parts matching the given criteria
func (s *Store) CountPartsByBrandModelYear(ctx context.Context, params models.SearchParams) (int, error) {
	query := `
		SELECT COUNT(*) FROM parts
		WHERE ($1 = '' OR brand ILIKE '%' || $1 || '%')
		AND ($2 = '' OR model ILIKE '%' || $2 || '%')
		AND ($3 = 0 OR year = $3)
		AND ($4 = '' OR category ILIKE '%' || $4 || '%')
	`

	var count int
	err := s.db.QueryRowContext(
		ctx,
		query,
		params.Brand,
		params.Model,
		params.Year,
		params.Category,
	).Scan(&count)

	if err != nil {
		return 0, fmt.Errorf("error counting parts: %w", err)
	}

	return count, nil
}

// SaveParts saves multiple parts to the database
func (s *Store) SaveParts(ctx context.Context, parts []models.Part) error {
	// Use a transaction for bulk insert
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		// If an error is returned, rollback the transaction
		if err != nil {
			tx.Rollback()
		}
	}()

	// Prepare the insert statement
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO parts (brand, model, year, name, category, description, image_url, source_url, price, currency, in_stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (brand, model, year, name) DO UPDATE SET
		description = EXCLUDED.description,
		image_url = EXCLUDED.image_url,
		source_url = EXCLUDED.source_url,
		price = EXCLUDED.price,
		currency = EXCLUDED.currency,
		in_stock = EXCLUDED.in_stock,
		updated_at = EXCLUDED.updated_at
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute for each part
	now := time.Now()
	for _, part := range parts {
		_, err := stmt.ExecContext(
			ctx,
			part.Brand,
			part.Model,
			part.Year,
			part.Name,
			part.Category,
			part.Description,
			part.ImageURL,
			part.SourceURL,
			part.Price,
			part.Currency,
			part.InStock,
			now,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert part: %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
