package products

import (
	"errors"
	"sync"
)

var ErrProductAlreadyExists = errors.New("product already exists")

type ProductRepository interface {
	Create(product Product) error
	GetByID(id int64) (Product, bool)
}

type repo struct {
	data map[int64]Product
	mu   sync.RWMutex
}

func NewProductRepository() ProductRepository {
	return &repo{
		data: make(map[int64]Product),
	}
}

func (r *repo) Create(product Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[product.ID]; ok {
		return ErrProductAlreadyExists
	}
	r.data[product.ID] = product
	return nil
}

func (r *repo) GetByID(id int64) (Product, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	product, ok := r.data[id]
	return product, ok
}
