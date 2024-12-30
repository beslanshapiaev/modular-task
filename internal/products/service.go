package products

import (
	"fmt"
	"modular-task/internal/eventbus"
	"modular-task/internal/users"
	"time"

	"math/rand"
)

type ProductService struct {
	repo        ProductRepository
	userService *users.UserService
	eventbus    *eventbus.EventBus
}

func NewProductService(repo ProductRepository, userService *users.UserService, eventBus *eventbus.EventBus) *ProductService {
	return &ProductService{
		repo:        repo,
		userService: userService,
		eventbus:    eventBus,
	}
}

func (s *ProductService) CreateProduct(userID int64, name string, price float64) (Product, error) {
	if _, ok := s.userService.GetUserByID(userID); !ok {
		return Product{}, fmt.Errorf("user with ID=%d not found", userID)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	product := Product{
		ID:     r.Int63n(1_000_000),
		UserID: userID,
		Name:   name,
		Price:  price,
	}
	err := s.repo.Create(product)
	if err != nil {
		return Product{}, err
	}
	s.eventbus.Publish(eventbus.Event{
		Type: eventbus.ProductCreated,
		Data: product,
	})
	return product, nil
}

func (s *ProductService) GetProductByID(id int64) (Product, bool) {
	return s.repo.GetByID(id)
}
