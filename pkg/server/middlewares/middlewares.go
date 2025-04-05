package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/wafi04/vazzuniversebackend/pkg/config"
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
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type JWTClaims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtSecretKey = []byte(config.LoadEnv("JWT_SECRET"))

func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Check for empty token
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check if token is expired
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func GenerateToken(user *UserData, hours int64) (string, error) {
	// Validate input
	if user == nil {
		return "", errors.New("user data cannot be nil")
	}

	if hours <= 0 {
		return "", errors.New("token duration must be positive")
	}

	claims := JWTClaims{
		UserID:   user.UserID,
		Email:    user.Email,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "wafiuddin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromGinContext(c *gin.Context) (*UserData, error) {
	user, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil, errors.New("user not found in context")
	}
	userInfo, ok := user.(*UserData)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}
	return userInfo, nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Try to get token from Authorization header first (Bearer token)
		authHeader := c.GetHeader("Authorization")
		var accessToken string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// Step 2: If no auth header, check cookie for access_token
			var err error
			accessToken, err = c.Cookie("access_token")
			if err != nil {
				accessToken = "" // Make sure it's empty if there's an error
			}
		}

		// Try to validate the access token if it exists
		if accessToken != "" {
			claims, err := ValidateToken(accessToken)
			if err == nil {
				user := &UserData{
					UserID:   claims.UserID,
					Email:    claims.Email,
					Username: claims.Username,
				}
				c.Set(string(UserContextKey), user)
				c.Next()
				return
			}
			log.Printf("Access token error: %v", err)
		}

		refreshToken, err := c.Cookie("refresh_token")
		if err == nil && refreshToken != "" {
			claims, err := ValidateToken(refreshToken)
			if err == nil {
				user := &UserData{
					UserID:   claims.UserID,
					Email:    claims.Email,
					Username: claims.Username,
				}

				newAccessToken, err := GenerateToken(user, 24)
				if err != nil {
					log.Printf("Failed to generate new access token: %v", err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"error":   "INTERNAL_SERVER_ERROR",
						"message": "Failed to generate new access token",
						"details": err.Error(),
					})
					return
				}

				// Set new access token cookie
				SetTokenCookie(c, "access_token", newAccessToken, 24*3600)

				c.Header("Authorization", "Bearer "+newAccessToken)

				// Set user in context
				c.Set(string(UserContextKey), user)
				c.Next()
				return
			}
			// Log refresh token error
			log.Printf("Refresh token error: %v", err)
		}

		// Step 4: No valid tokens found
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "UNAUTHORIZED",
			"message": "Authentication required",
			"details": "No valid tokens found",
		})
	}
}

func SetTokenCookie(c *gin.Context, name, token string, duration int) {
	// Secure and HttpOnly flags should be true in production
	secure := false
	if config.LoadEnv("ENV") == "production" {
		secure = true
	}

	c.SetCookie(
		name,
		token,
		duration,
		"/",
		"localhost",
		secure,
		true,
	)

	log.Printf("Cookie %s set with duration %d seconds", name, duration)
}

func ClearTokens(c *gin.Context) {
	// Clear refresh token
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	// Clear access token
	c.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	// Clear Authorization header
	c.Header("Authorization", "")

	log.Printf("All authentication tokens cleared")
}

func ResponseTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		c.Header("X-Response-Time", duration.String())

		// Optional: Log response time for monitoring
		log.Printf("[%s] %s - %s", c.Request.Method, c.Request.URL.Path, duration)
	}
}
