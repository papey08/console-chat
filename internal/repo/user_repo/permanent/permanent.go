package permanent

import (
	"console-chat/internal/model"
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	// addUserQuery is a query to insert user into database
	addUserQuery = `
		INSERT INTO users (nickname, hashed_password)
		VALUES ($1, $2);`

	// getUserQuery is a query to select user from the database
	getUserQuery = `
		SELECT * FROM users
		WHERE nickname = $1;`
)

// duplicateCode is a code of pgconn.PgError when primary key is duplicated
const duplicateCode = "23505"

// permanentRepo is a permanent storage of all users
type PermanentRepo struct {
	pgx.Conn
}

func (r *PermanentRepo) InsertUser(ctx context.Context, u model.User) (model.User, error) {
	_, err := r.Exec(ctx, addUserQuery, u.Nickname, u.HashedPassword)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == duplicateCode {
			return model.User{}, model.UserAlreadyExists
		} else {
			// debug info
			log.Println(err.Error())
			return model.User{}, model.UserRepoError
		}
	}
	return u, nil
}

func (r *PermanentRepo) SelectUser(ctx context.Context, nickname string) (model.User, error) {
	var usr model.User
	row := r.QueryRow(ctx, getUserQuery, nickname)
	if err := row.Scan(&usr.Nickname, &usr.HashedPassword); err == pgx.ErrNoRows {
		return model.User{}, model.UserNotFound
	} else if err != nil {
		return model.User{}, model.UserRepoError
	} else {
		return usr, nil
	}
}
