package middlewares

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserData struct {
	UserID    string     `json:"userId" db:"user_id"`
	Fullname  *string    `json:"fullName"  db:"full_name"`
	Username  string     `json:"username"  db:"username"`
	Email     string     `json:"email"  db:"email"`
	Password  *string    `json:"password,omitempty"  db:"password"`
	Role      string     `json:"role"  db:"role"`
	IsDeleted bool       `json:"isDeleted"  db:"is_deleted"`
	Balance   int        `json:"balance" db:"balance"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	SessionID string     `json:"sessionId" db:"session_id"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
	jwt.RegisteredClaims
}
