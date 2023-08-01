package cache

import (
	"console-chat/internal/model"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

// expiration is how long new user would stay in cache after registration
const expiration = time.Minute * 30

type cachedUser struct {
	Nickname       string `json:"nickname"`
	HashedPassword string `json:"hashed_password"`
}

func usrToCashedUsr(u model.User) cachedUser {
	return cachedUser{
		Nickname:       u.Nickname,
		HashedPassword: u.HashedPassword,
	}
}

func cachedUsrToUsr(u cachedUser) model.User {
	return model.User{
		Nickname:       u.Nickname,
		HashedPassword: u.HashedPassword,
	}
}

// cacheRepo is a temporary storage for new registered users to make their
// signing in faster
type CacheRepo struct {
	redis.Client
}

/*
func New(rc *redis.Client) userrepo.Cache {
	return &cacheRepo{
		Client: *rc,
	}
} */

func (c *CacheRepo) SetUserByKey(ctx context.Context, key string, u model.User) (model.User, error) {
	cu := usrToCashedUsr(u)
	data, _ := json.Marshal(cu)
	if err := c.Set(ctx, u.Nickname, data, expiration).Err(); err != nil {
		return model.User{}, model.UserRepoError
	}
	return u, nil
}

func (c *CacheRepo) GetUserByKey(ctx context.Context, key string) (model.User, error) {
	recievedData, err := c.Get(ctx, key).Result()
	if err == redis.Nil {
		return model.User{}, model.UserNotFound
	} else if err != nil {
		return model.User{}, model.UserRepoError
	}

	var recievedUser cachedUser
	if err := json.Unmarshal([]byte(recievedData), &recievedUser); err != nil {
		return model.User{}, model.UserRepoError
	} else {
		return cachedUsrToUsr(recievedUser), nil
	}
}
