package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// user struct to represent a indvidual user
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// user model that wraps a db connection pool
type UserModel struct {
	DB *pgxpool.Pool
}

// create a new user record in the db
func (u *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	statement := `INSERT INTO users (name, email, hashed_password, created) VALUES ($1, $2, $3, NOW() AT TIME ZONE 'UTC');`

	args := []any{
		name,
		email,
		hashedPassword,
	}

	_, err = u.DB.Exec(context.Background(), statement, args...)
	if err != nil {
		var postgreqlError *pgconn.PgError

		if errors.As(err, &postgreqlError) {
			// if error is code 23505 then we have duplicate email so return
			// our ErrDuplicateEmail error defined in errors.go.
			if postgreqlError.Code == "23505" && strings.Contains(postgreqlError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

// Verifys whether a user exists with the provided email and password.
// This will return the relevant user ID if they do.
func (u *UserModel) Authenticate(email, password string) (int, error) {

	// Check if user exists
	var id int
	var hashedPassword []byte
	statement := "select id, hashed_password from users where email = $1"
	err := u.DB.QueryRow(context.Background(), statement, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// Check if the password match the stored hashed password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

// Used to check if a user exists with a specific ID
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
