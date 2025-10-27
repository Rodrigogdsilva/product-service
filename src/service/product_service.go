package service

import (
	"context"
	"fmt"
	"product-service/src/domain"
	"product-service/src/repository"
	"time"

	"github.com/google/uuid"
)

type ProductService interface {
	Create(ctx context.Context, name, description string, price float64, stock int) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) Create(ctx context.Context, name, description string, price float64, stock int) error {

	if name == "" || description == "" {
		return fmt.Errorf("Error creating product: %w", domain.ErrParametersMissing)
	}
	if price <= 0 {
		return fmt.Errorf("Error creating product: %w", domain.ErrInvalidPrice)
	}
	if stock < 0 {
		return fmt.Errorf("Error creating product: %w", domain.ErrInvalidStock)
	}

	product := &domain.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	return s.productRepository.Create(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {

	if id == uuid.Nil {
		return nil, fmt.Errorf("Error when searching for product by ID: %w", domain.ErrInvalidID)
	}

	return s.productRepository.GetProductByID(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.productRepository.ListProducts(ctx)
}

func (s *productService) ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error {

	if id == uuid.Nil {
		return fmt.Errorf("Error when reducing stock: %w", domain.ErrInvalidID)
	}

	if quantity <= 0 {
		return fmt.Errorf("Error when reducing stock: %w", domain.ErrInvalidQuantity)
	}

	return s.productRepository.ReduceStock(ctx, id, quantity)
}

func (s *productService) Update(ctx context.Context, product *domain.Product) error {

	if product.Name == "" || product.Description == "" {
		return fmt.Errorf("Error updating product: %w", domain.ErrParametersMissing)
	}
	if product.Price <= 0 {
		return fmt.Errorf("Error updating product: %w", domain.ErrInvalidPrice)
	}
	if product.Stock < 0 {
		return fmt.Errorf("Error updating product: %w", domain.ErrInvalidStock)
	}

	product.UpdatedAt = time.Now().UTC()

	return s.productRepository.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {

	if id == uuid.Nil {
		return fmt.Errorf("Error when deleting product: %w", domain.ErrInvalidID)
	}

	return s.productRepository.Delete(ctx, id)
}
