package repository

import (
	"context"
	"fmt"
	"product-service/src/domain"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetProductByID(ctx context.Context, id ulid.ULID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	ReduceStock(ctx context.Context, id ulid.ULID, quantity int) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id ulid.ULID) error
}

type postgresProductRepository struct {
	db *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) ProductRepository {
	return &postgresProductRepository{db: db}
}

func (r *postgresProductRepository) Create(ctx context.Context, product *domain.Product) error {

	if product.Name == "" || product.Description == "" || product.Price == 0 || product.Stock == 0 {
		return fmt.Errorf("Error creating product: %w", domain.ErrParametersMissing)
	}

	if product.Price < 0 {
		return fmt.Errorf("Error creating product: %w", domain.ErrInvalidPrice)
	}

	if product.Stock < 0 {
		return fmt.Errorf("Error creating product: %w", domain.ErrInvalidStock)
	}

	query := `INSERT INTO products (id, name, description, price, stock) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, product.Name, product.Description, product.Price)
	if err != nil {
		return fmt.Errorf("Error creating product: %w", domain.ErrFailedCreatingProduct)
	}
	return nil
}

func (r *postgresProductRepository) GetProductByID(ctx context.Context, id ulid.ULID) (*domain.Product, error) {

	if id == (ulid.ULID{}) {
		return nil, fmt.Errorf("Error when searching for product by ID: %w", domain.ErrInvalidID)
	}

	query := `SELECT id, name, description, price, stock FROM products WHERE id = $1` // Assuming ID is a ULID, this query might need adjustment depending on how ULID is stored in DB
	product := &domain.Product{}
	err := r.db.QueryRow(ctx, query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
	if err != nil {
		return nil, fmt.Errorf("Error when searching for product by ID: %w", domain.ErrProductNotFound)
	}
	return product, nil
}

func (r *postgresProductRepository) ListProducts(ctx context.Context) ([]*domain.Product, error) {

	query := `SELECT id, name, description, price, stock FROM products`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Error when searching for all products: %w", domain.ErrNotFoundProducts)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock)
		if err != nil {
			return nil, fmt.Errorf("Error scanning product row: %w", domain.ErrProductNotFound)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *postgresProductRepository) ReduceStock(ctx context.Context, id ulid.ULID, quantity int) error {
	if id == (ulid.ULID{}) {
		return fmt.Errorf("Error when reducing stock: %w", domain.ErrInvalidID)
	}
	if quantity <= 0 {
		return fmt.Errorf("Error when reducing stock: %w", domain.ErrInvalidQuantity)
	}

	query := `UPDATE products SET stock = stock - $1 WHERE id = $2`
	err := r.db.QueryRow(ctx, query, quantity, id)
	if err != nil {
		return fmt.Errorf("Error when reducing stock: %w", domain.ErrToReduceStock)
	}
	return nil
}

func (r *postgresProductRepository) Update(ctx context.Context, product *domain.Product) error {

	if product.ID == (ulid.ULID{}) { // Check if ULID is zero value
		return fmt.Errorf("Error when updating product: %w", domain.ErrInvalidID)
	}

	query := `UPDATE products SET name = $1, description = $2, price = $3, stock = $4 WHERE id = $5`
	err := r.db.QueryRow(ctx, query, product.Name, product.Description, product.Price, product.Stock, product.ID)
	if err != nil {
		return fmt.Errorf("Error when updating product: %w", domain.ErrToUpdateProduct)
	}
	return nil
}

func (r *postgresProductRepository) Delete(ctx context.Context, id ulid.ULID) error {

	if id == (ulid.ULID{}) {
		return fmt.Errorf("Error when deleting product: %w", domain.ErrInvalidID)
	}

	query := `DELETE FROM products WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id)
	if err != nil {
		return fmt.Errorf("Error when deleting product: %w", domain.ErrToDeletegProduct)
	}
	return nil
}
