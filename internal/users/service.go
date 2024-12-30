package users

import (
	"math/rand"
	"modular-task/internal/eventbus"
	"time"
)

type UserService struct {
	repo     UserRepository
	eventBus *eventbus.EventBus
}

func NewUserService(repo UserRepository, eventBus *eventbus.EventBus) *UserService {
	return &UserService{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *UserService) CreateUser(name, email string) (User, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	user := User{
		ID:    r.Int63n(1_000_000),
		Name:  name,
		Email: email,
	}
	err := s.repo.Create(user)
	if err != nil {
		return User{}, err
	}

	s.eventBus.Publish(eventbus.Event{
		Type: eventbus.EventUserCreated,
		Data: user,
	})
	return user, nil
}

func (s *UserService) GetUserByID(id int64) (User, bool) {
	return s.repo.GetByID(id)
}
