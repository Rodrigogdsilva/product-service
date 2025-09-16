package service

import (
	"context"
	"math/rand"
	"product-service/src/domain"
	"product-service/src/repository"
	"time"

	"github.com/oklog/ulid/v2"
)

type ProductService interface {
	Create(ctx context.Context, name, description string, price float64, stock int) error
	GetProductByID(ctx context.Context, id ulid.ULID) (*domain.Product, error)
	ListProducts(ctx context.Context) ([]*domain.Product, error)
	ReduceStock(ctx context.Context, id ulid.ULID, quantity int) error
	Update(ctx context.Context, userID ulid.ULID, name, description string, price float64, stock int) error
	Delete(ctx context.Context, id ulid.ULID) error
}

type productService struct {
	productRepository repository.ProductRepository
}

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository: productRepository}
}

func (s *productService) Create(ctx context.Context, name, description string, price float64, stock int) error {

	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	newID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

	product := &domain.Product{
		ID:          newID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	return s.productRepository.Create(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, id ulid.ULID) (*domain.Product, error) {
	return s.productRepository.GetProductByID(ctx, id)
}

func (s *productService) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	return s.productRepository.ListProducts(ctx)
}

func (s *productService) ReduceStock(ctx context.Context, id ulid.ULID, quantity int) error {
	return s.productRepository.ReduceStock(ctx, id, quantity)
}

func (s *productService) Update(ctx context.Context, userID ulid.ULID, name, description string, price float64, stock int) error {

	product := &domain.Product{
		ID:          userID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
	}

	return s.productRepository.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id ulid.ULID) error {
	return s.productRepository.Delete(ctx, id)
}
