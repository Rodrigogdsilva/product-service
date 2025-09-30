package service

import (
	"context"
	"product-service/src/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) Create(ctx context.Context, name, description string, price float64, stock int) error {
	args := m.Called(ctx, name, description, price, stock)
	return args.Error(0)
}

func (m *ProductServiceMock) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if product, ok := args.Get(0).(*domain.Product); ok {
		return product, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) ListProducts(ctx context.Context) ([]*domain.Product, error) {
	args := m.Called(ctx)
	if products, ok := args.Get(0).([]*domain.Product); ok {
		return products, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ProductServiceMock) ReduceStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *ProductServiceMock) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductServiceMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
