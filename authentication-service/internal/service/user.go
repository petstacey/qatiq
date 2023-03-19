package service

import (
	"context"
	"errors"
	"time"

	"github.com/pso-dev/qatiq/backend/authentication-service/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Active   int64  `json:"active"`
}

func ValidateID(v *validator.Validator, id int64) {
	v.Check(id != 0, "id", "must be provided")
	v.Check(id > 0, "id", "cannot be a negative number")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 characters long")
}

func ValidayeUser(v *validator.Validator, user User) {
	ValidateID(v, user.ID)
	ValidateEmail(v, user.Email)
	ValidatePassword(v, user.Password)
}

func (s *Service) GetAll() ([]*User, error) {
	query := `SELECT id, email, password_hash, active FROM users`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Active,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (s *Service) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, password_hash, active FROM users WHERE email=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user User
	err := s.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Active,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Service) PasswordMatches(hash, plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
