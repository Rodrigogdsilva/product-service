package service

import (
	"context"
	"product-service/src/domain"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/mock"
)

type ProductServiceMock struct {
	mock.Mock
}

func (m *ProductServiceMock) Create(ctx context.Context, name, description string, price float64, stock int) error {
	args := m.Called(ctx, name, description, price, stock)
	return args.Error(0)
}

func (m *ProductServiceMock) GetProductByID(ctx context.Context, id ulid.ULID) (*domain.Product, error) {
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

func (m *ProductServiceMock) ReduceStock(ctx context.Context, id ulid.ULID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *ProductServiceMock) Update(ctx context.Context, userID ulid.ULID, name, description string, price float64, stock int) error {
	args := m.Called(ctx, userID, name, description, price, stock)
	return args.Error(0)
}

func (m *ProductServiceMock) Delete(ctx context.Context, id ulid.ULID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
