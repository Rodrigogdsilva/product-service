package seeder

import (
	"context"
	"product-service/src/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestSeeder struct {
	db *pgxpool.Pool
}

func NewTestSeeder(db *pgxpool.Pool) *TestSeeder {
	return &TestSeeder{db: db}
}

func (s *TestSeeder) InsertProduct(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (id, name, description, price, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := s.db.Exec(ctx, query, product.ID, product.Name, product.Description, product.Price, product.Stock, product.CreatedAt, product.UpdatedAt)
	return err
}
