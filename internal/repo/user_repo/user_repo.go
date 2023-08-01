package userrepo

import (
	"console-chat/internal/app"
	"console-chat/internal/model"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
)

type permanent interface {
	InsertUser(ctx context.Context, u model.User) (model.User, error)
	SelectUser(ctx context.Context, nickname string) (model.User, error)
}

type cache interface {
	SetUserByKey(ctx context.Context, key string, u model.User) (model.User, error)
	GetUserByKey(ctx context.Context, key string) (model.User, error)
}

type Repo struct {
	permanent
	cache
}

func New(conn *pgx.Conn, rc *redis.Client) app.UserRepo {
	return &Repo{
		permanent: newPermanent(conn),
		cache:     newCache(rc),
	}
}

func (r *Repo) AddUser(ctx context.Context, u model.User) (model.User, error) {
	usr, err := r.InsertUser(ctx, u) // add user to permanent db
	if err != nil {
		return model.User{}, err
	}
	_, err = r.SetUserByKey(ctx, u.Nickname, u) // add user to cache
	if err != nil {
		return model.User{}, err
	}
	return usr, nil
}

func (r *Repo) GetUser(ctx context.Context, nickname string) (model.User, error) {
	if usr, err := r.GetUserByKey(ctx, nickname); err == model.UserRepoError { // case when something wrong with cache
		return model.User{}, err
	} else if err == nil { // case when user was found in cache
		return usr, nil
	}

	// case when user is not in cache
	if usr, err := r.SelectUser(ctx, nickname); err != nil { // case when usr not in cache and not in db
		return model.User{}, err
	} else { // case when user in db but not in cache
		if _, err := r.SetUserByKey(ctx, usr.Nickname, usr); err != nil {
			return model.User{}, err
		}
		return usr, nil
	}

}