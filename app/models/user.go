package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"project/internal/platform/crypto"
	"project/internal/platform/database"
)

// User represents a user entity
type User struct {
	Entity
	Email     string `db:"email"`
	Password  string `db:"password"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	IsActive  bool   `db:"is_active"`
}

// UserStore represents a user data store
type UserStore struct {
	Model
}

// NewUserStore creates a new UserStore instance
func NewUserStore(db *database.Database) *UserStore {
	return &UserStore{
		Model: *NewModel(db),
	}
}

// FindByID finds a user by ID
func (s *UserStore) FindByID(ctx context.Context, id string) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`

	qCtx := s.GetQueryableContext(ctx)
	user := &User{}
	row := qCtx.Queryable().QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// FindByEmail finds a user by email
func (s *UserStore) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`

	qCtx := s.GetQueryableContext(ctx)
	user := &User{}
	row := qCtx.Queryable().QueryRowContext(ctx, query, email)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

// FindAll finds all users
func (s *UserStore) FindAll(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, email, password, first_name, last_name, is_active, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`

	qCtx := s.GetQueryableContext(ctx)
	rows, err := qCtx.Queryable().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Create creates a new user
func (s *UserStore) Create(ctx context.Context, user *User) error {
	// Hash the password
	hashedPassword, err := crypto.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Set the hashed password
	user.Password = hashedPassword

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, password, first_name, last_name, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	qCtx := s.GetQueryableContext(ctx)
	_, err = qCtx.Queryable().ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

// Update updates an existing user
func (s *UserStore) Update(ctx context.Context, user *User) error {
	// Update timestamp
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = ?, first_name = ?, last_name = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	qCtx := s.GetQueryableContext(ctx)
	_, err := qCtx.Queryable().ExecContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	return err
}

// UpdatePassword updates a user's password
func (s *UserStore) UpdatePassword(ctx context.Context, id, password string) error {
	// Hash the password
	hashedPassword, err := crypto.HashPassword(password)
	if err != nil {
		return err
	}

	// Update the password
	query := `
		UPDATE users
		SET password = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	qCtx := s.GetQueryableContext(ctx)
	_, err = qCtx.Queryable().ExecContext(ctx, query,
		hashedPassword,
		time.Now(),
		id,
	)

	return err
}

// Delete soft deletes a user
func (s *UserStore) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE users
		SET deleted_at = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	qCtx := s.GetQueryableContext(ctx)
	_, err := qCtx.Queryable().ExecContext(ctx, query,
		time.Now(),
		time.Now(),
		id,
	)

	return err
}

// Authenticate authenticates a user with email and password
func (s *UserStore) Authenticate(ctx context.Context, email, password string) (*User, error) {
	// Find the user by email
	user, err := s.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Check if the user exists
	if user == nil {
		return nil, nil
	}

	// Check if the password is correct
	if !crypto.CheckPasswordHash(password, user.Password) {
		return nil, nil
	}

	return user, nil
}
