package users

import (
	"errors"
	"sync"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository interface {
	Create(user User) error
	GetByID(id int64) (User, bool)
}

type repo struct {
	data map[int64]User
	mu   sync.Mutex
}

func NewUserRepository() UserRepository {
	return &repo{
		data: make(map[int64]User),
	}
}

func (r *repo) Create(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[user.ID]; ok {
		return ErrUserAlreadyExists
	}
	r.data[user.ID] = user
	return nil
}

func (r *repo) GetByID(id int64) (User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.data[id]
	return user, ok
}
