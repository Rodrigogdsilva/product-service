package repository

import (
	"context"
	"errors"
	"fmt"
	"product-service/src/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type postgresProductRepository struct {
	db *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *domain.Product) error {

	query := `INSERT INTO products (id, name, description, price, stock, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, product.ID, product.Name, product.Description, product.Price, product.Stock, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		return fmt.Errorf("Error creating product: %w", domain.ErrFailedCreatingProduct)
	}
	return nil
}

func (r *postgresProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {

	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id = $1`
	product := &domain.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("Error when searching for product by ID: %w", domain.ErrProductNotFound)
		}
		return nil, fmt.Errorf("Error when searching for product by ID: %w", err)
	}
	return product, nil
}

func (r *postgresProductRepository) ListProducts(ctx context.Context) ([]*domain.Product, error) {

	query := `SELECT id, name, description, price, stock, created_at, updated_at FROM products`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Error when searching for all products: %w", domain.ErrNotFoundProducts)
	}
	defer rows.Close()

	products := make([]*domain.Product, 0)
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning product row: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *postgresProductRepository) ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error {

	query := `UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, quantity, time.Now(), id)
	if err != nil {
		return fmt.Errorf("Error when reducing stock: %w", domain.ErrToReduceStock)
	}
	return nil
}

func (r *postgresProductRepository) Update(ctx context.Context, product *domain.Product) error {

	query := `UPDATE products SET name = $1, description = $2, price = $3, stock = $4, updated_at = $5 WHERE id = $6`
	_, err := r.db.Exec(ctx, query, product.Name, product.Description, product.Price, product.Stock, time.Now(), product.ID)
	if err != nil {
		return fmt.Errorf("Error when updating product: %w", domain.ErrToUpdateProduct)
	}
	return nil
}

func (r *postgresProductRepository) Delete(ctx context.Context, id uuid.UUID) error {

	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Error when deleting product: %w", domain.ErrToDeletegProduct)
	}
	return nil
}
