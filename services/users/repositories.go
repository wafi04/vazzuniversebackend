package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/generate"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
)

type UserRepositories struct {
	MainDB    *sqlx.DB
	ReplicaDB *sqlx.DB
}

func NewUserRepositories(MainDB, ReplicaDB *sqlx.DB) *UserRepositories {
	return &UserRepositories{
		MainDB:    MainDB,
		ReplicaDB: ReplicaDB,
	}
}

// Create inserts a new user into the database
func (r *UserRepositories) Create(ctx context.Context, user CreateUser) (*UserData, *response.ResponseError) {
	// Generate UUID for user_id
	userID := generate.GenerateRandomID(&generate.IDOpts{
		Amount: 10,
	})

	hashedPassword, err := HashPassword(*user.Password)
	if err != nil {
		return nil, ErrInternalServerError(err)
	}

	now := time.Now()

	query := `
    INSERT INTO users (user_id,full_name, username, email, password, role, balance, is_deleted, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9 ,$10)
    RETURNING user_id, full_name, username, email, role, balance, is_deleted, created_at, updated_at
    `

	// Create user data
	userData := UserData{
		UserID:    userID,
		Fullname:  user.FullName,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Balance:   user.Balance,
		IsDeleted: false,
		CreatedAt: now,
		UpdatedAt: &now,
	}

	// Execute query
	var result UserData
	err = r.MainDB.QueryRowxContext(ctx, query,
		userData.UserID,
		userData.Fullname,
		userData.Username,
		userData.Email,
		hashedPassword,
		userData.Role,
		userData.Balance,
		userData.IsDeleted,
		userData.CreatedAt,
		userData.UpdatedAt,
	).StructScan(&result)

	if err != nil {
		// Better error handling for database errors
		return nil, ErrInternalServerError(err)
	}

	return &result, nil
}

func (r *UserRepositories) GetUserByID(ctx context.Context, userID string) (*UserData, error) {
	query := `
    SELECT user_id,full_name, username, email, role, balance, is_deleted, created_at, updated_at
    FROM users
    WHERE user_id = $1 AND is_deleted = false
    `

	var userData UserData
	err := r.ReplicaDB.GetContext(ctx, &userData, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &userData, nil
}

func (r *UserRepositories) GetUserByUsername(ctx context.Context, username string) (*UserData, error) {
	query := `
    SELECT user_id,full_name, username, email, role, balance, is_deleted, created_at, updated_at
    FROM users
    WHERE username = $1 AND is_deleted = false
    `

	var userData UserData
	err := r.ReplicaDB.GetContext(ctx, &userData, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &userData, nil
}

func (r *UserRepositories) GetUserByEmail(ctx context.Context, email string) (*UserData, error) {
	query := `
    SELECT user_id,full_name, username, email, role, balance, is_deleted, created_at, updated_at
    FROM users
    WHERE email = $1 AND is_deleted = false
    `

	var userData UserData
	err := r.ReplicaDB.GetContext(ctx, &userData, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return &userData, nil
}
