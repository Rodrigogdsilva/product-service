package stubs

import (
	"math/rand"
	"product-service/src/domain"
	"time"

	"github.com/jaswdr/faker"
	"github.com/oklog/ulid/v2"
)

type ProductStub struct {
	product *domain.Product
}

func NewProductStub() *ProductStub {
	f := faker.New()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)

	return &ProductStub{
		product: &domain.Product{
			ID:          ulid.MustNew(ulid.Timestamp(time.Now()), entropy),
			Name:        f.Person().Name(),
			Description: f.Lorem().Sentence(10),
			Price:       f.Float64(2, 10, 1000),
			Stock:       f.IntBetween(1, 100),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}

func (s *ProductStub) WithID(id ulid.ULID) *ProductStub {
	s.product.ID = id
	return s
}

func (s *ProductStub) WithName(name string) *ProductStub {
	s.product.Name = name
	return s
}

func (s *ProductStub) WithPrice(price float64) *ProductStub {
	s.product.Price = price
	return s
}

func (s *ProductStub) Get() *domain.Product {
	return s.product
}
