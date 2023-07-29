package userrepo

import (
	"console-chat/internal/model"
	"context"
	"sync"
)

type Repo struct {
	usrs map[string]model.User
	mu   *sync.Mutex
}

func New() *Repo {
	return &Repo{
		usrs: make(map[string]model.User),
		mu:   new(sync.Mutex),
	}
}

func (r *Repo) AddUser(ctx context.Context, u model.User) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-ctx.Done():
		return model.User{}, model.UserRepoError
	default:
		if _, ok := r.usrs[u.Nickame]; ok {
			return model.User{}, model.UserAlreadyExists
		} else {
			r.usrs[u.Nickame] = u
			return u, nil
		}
	}
}

func (r *Repo) GetUser(ctx context.Context, nickname string) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	select {
	case <-ctx.Done():
		return model.User{}, model.UserRepoError
	default:
		if u, ok := r.usrs[nickname]; ok {
			return u, nil
		} else {
			return model.User{}, model.UserNotFound
		}
	}
}
