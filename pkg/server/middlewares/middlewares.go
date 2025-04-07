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
	"github.com/wafi04/vazzuniversebackend/pkg/constants"
	"github.com/wafi04/vazzuniversebackend/pkg/utils/response"
	"github.com/wafi04/vazzuniversebackend/services/auth/sessions"
)

func ValidateToken(tokenString string) (*JWTClaims, error) {
	config.LoadEnv("JWT_SECRET")
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return constants.JWT_SECRET, nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

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
		UserID:    user.UserID,
		Email:     user.Email,
		Username:  user.Username,
		SessionID: user.SessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(hours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    config.LoadEnv("AUTHOR"),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(constants.JWT_SECRET)
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
func AuthMiddleware(sessionService *sessions.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var accessToken string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			accessToken = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			var err error
			accessToken, err = c.Cookie("auth_token")
			if err != nil {
				accessToken = ""
				response.NewResponseError(http.StatusUnauthorized, "USER UNATHORIZED", "UNAUTHORIZED")
			} else {
				response.NewResponseError(http.StatusUnauthorized, "USER UNATHORIZED", "UNAUTHORIZED")
			}
		}

		log.Printf("%s", accessToken)
		if accessToken != "" {
			claims, err := ValidateToken(accessToken)
			if err == nil {
				// Validate session using sessionID from JWT claims
				sessionData, err := sessionService.GetBySessionID(c.Request.Context(), claims.SessionID)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"success": false,
						"error":   "UNAUTHORIZED",
						"message": "Session is invalid or expired",
						"details": "Invalid session",
					})
					return
				}

				// Check if session is still active (not invalidated)
				if !sessionData.IsAccess {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"success": false,
						"error":   "UNAUTHORIZED",
						"message": "Session has been invalidated",
						"details": "Inactive session",
					})
					return
				}

				err = sessionService.UpdateLastActivity(c.Request.Context(), claims.SessionID)
				if err != nil {
					response.NewResponseError(http.StatusBadRequest, "Sessions Invalidate", "Please Login Frist!")
				}

				user := &UserData{
					UserID:    claims.UserID,
					Email:     claims.Email,
					Username:  claims.Username,
					SessionID: claims.SessionID,
				}
				c.Set(string(UserContextKey), user)
				c.Next()
				return
			}
		}

		// No valid tokens found
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "UNAUTHORIZED",
			"message": "Authentication required",
			"details": "No valid tokens found",
		})
	}
}
func SetTokenCookie(c *gin.Context, name, token string, duration int) {
	// SameSite attribute penting untuk localhost
	c.SetSameSite(http.SameSiteNoneMode)

	c.SetCookie(
		name,
		token,
		duration,
		"/",
		"",    // Kosongkan domain untuk localhost
		false, // Secure harus false untuk http
		false, // Gunakan false sementara untuk debugging (biasanya true)
	)

	// Tambahkan juga sebagai header untuk debug
	c.Header("Set-Cookie", name+"="+token+"; Path=/; Max-Age="+fmt.Sprint(duration))

}
func ClearTokens(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"",
		false,
		false,
	)

	// Clear Authorization header
	c.Header("Authorization", "")
}

func ResponseTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		c.Header("X-Response-Time", duration.String())
	}
}
