package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/wafi04/vazzuniversebackend/pkg/utils"
)

type UserRepositories struct {
    DB *sqlx.DB
}

func NewUserRepositories(DB *sqlx.DB) *UserRepositories {
    return &UserRepositories{
        DB: DB,
    }
}

// Create inserts a new user into the database
func (r *UserRepositories) Create(ctx context.Context, user CreateUSer) (*UserData, *utils.ResponseError) {
    // Generate UUID untuk user_id
    userID := uuid.New().String()

    // Hash password sebelum disimpan
    hashedPassword, err := HashPassword(*user.Password)
    if err != nil {
        return nil, ErrInvalidCredentialsError()
    }

    now := time.Now()

    query := `
    INSERT INTO users (user_id, username, email, password, role, balance, is_deleted, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING user_id, username, email, role, balance, is_deleted, created_at, updated_at
    `

    // Buat user data
    userData := UserData{
        UserID:    userID,
        Username:  user.Username,
        Email:     user.Email,
        Role:      user.Role,
        Balance:   user.Balance,
        IsDeleted: false,
        CreatedAt: now,
        UpdatedAt: &now,
    }

    // Eksekusi query
    row := r.DB.QueryRowxContext(ctx, query,
        userData.UserID,
        userData.Username,
        userData.Email,
        hashedPassword,
        userData.Role,
        userData.Balance,
        userData.IsDeleted,
        userData.CreatedAt,
        userData.UpdatedAt,
    )

    var result UserData
    err = row.StructScan(&result)
    if err != nil {
        return nil, ErrEmailTakenError(user.Email)
    }

    return &result, nil
}

func (r *UserRepositories) GetUserByID(ctx context.Context, userID string) (*UserData, error) {
    query := `
    SELECT user_id, username, email, role, balance, is_deleted, created_at, updated_at
    FROM users
    WHERE user_id = $1 AND is_deleted = false
    `

    var userData UserData
    err := r.DB.GetContext(ctx, &userData, query, userID)
    if err != nil {
        return nil, err
    }

    return &userData, nil
}


func (r *UserRepositories) GetUserByUsername(ctx context.Context, username string) (*UserData, error) {
    query := `
    SELECT user_id, username, email, role, balance, is_deleted, created_at, updated_at
    FROM users
    WHERE username = $1 AND is_deleted = false
    `

    var userData UserData
    err := r.DB.GetContext(ctx, &userData, query, username)
    if err != nil {
        return nil, err
    }

    return &userData, nil
}
func (r *UserRepositories) GetUserByEmail(ctx context.Context, email string) (*UserData, error) {
    query := `
    SELECT user_id, username, email, role, balance, is_deleted, created_at, updated_at
    FROM users
    WHERE email = $1 AND is_deleted = false
    `

    var userData UserData
    err := r.DB.GetContext(ctx, &userData, query, email)
    if err != nil {
        return nil, err
    }

    return &userData, nil
}
