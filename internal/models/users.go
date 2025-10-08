package models

import (
	"time"

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
	_, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	return nil
}

// Verifys whether a user exists with the provided email and password.
// This will return the relevant user ID if they do.
func (u *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Used to check if a user exists with a specific ID
func (u *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
