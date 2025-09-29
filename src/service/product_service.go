package service

import (
	"context"
	"product-service/src/domain"
	"product-service/src/repository"

	"github.com/google/uuid"
)

type ProductService interface {
	Create(ctx context.Context, name, description string, price float64, stock int) error
	GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error
	Update(ctx context.Context, userID uuid.UUID, name, description string, price float64, stock int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) Create(ctx context.Context, name, description string, price float64, stock int) error {
	product := &domain.Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	return s.productRepository.Create(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return s.productRepository.GetProductByID(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.productRepository.ListProducts(ctx)
}

func (s *productService) ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return s.productRepository.ReduceStock(ctx, id, quantity)
}

func (s *productService) Update(ctx context.Context, userID uuid.UUID, name, description string, price float64, stock int) error {

	product := &domain.Product{
		ID:          userID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	return s.productRepository.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.productRepository.Delete(ctx, id)
}
