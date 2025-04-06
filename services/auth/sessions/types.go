package sessions

import "time"

var (
	ErrSessionsInvalid  = "INVALID_SESSIONS"
	ErrSessionsNotFound = "SESSIONS_NOT_FOUND"
	ErrSessionsExpired  = "SESSIONS_EXPIRED"
)

type SessionsData struct {
	SessionID    string     `json:"sessionId"      db:"session_id"`
	UserID       string     `json:"userId"         db:"user_id"`
	AccessToken  string     `json:"accessToken"    db:"access_token"`
	IPAddress    string     `json:"ipAddress"      db:"ip_address"`
	UserAgent    string     `json:"userAgent"      db:"user_agent"`
	DeviceInfo   string     `json:"deviceInfo"     db:"device_info"`
	LastActivity *time.Time `json:"lastActivity"   db:"last_activity"`
	ExpiresAt    *time.Time `json:"expiresAt"      db:"expires_at"`
	IsAccess     bool       `json:"isAccess"       db:"is_access"`
	CreatedAt    time.Time  `json:"createdAt"      db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt"      db:"updated_at"`
}

type CreateSession struct {
	SessionID    string    `json:"sessionId"    db:"session_id"`
	UserID       string    `json:"userId"       db:"user_id"`
	AccessToken  string    `json:"accessToken"  db:"access_token"`
	IPAddress    string    `json:"ipAddress"    db:"ip_address"`
	UserAgent    string    `json:"userAgent"    db:"user_agent"`
	DeviceInfo   string    `json:"deviceInfo"   db:"device_info"`
	ExpiresAt    time.Time `json:"expiresAt"    db:"expires_at"`
	LastActivity time.Time `json:"lastActivity"  db:"last_activity"`
	IsAccess     bool      `json:"isAccess"     db:"is_access"`
}

type CheckValidaty struct {
	SessionID string `json:"sessionId"  db:"session_id"`
	UserID    string `json:"userId"     db:"user_id"`
}

type UpdateSessionActivity struct {
	SessionID    string    `json:"sessionId"    db:"session_id"`
	LastActivity time.Time `json:"lastActivity" db:"last_activity"`
}

type SessionFilter struct {
	UserID     string `json:"userId"     db:"user_id"`
	IPAddress  string `json:"ipAddress"  db:"ip_address"`
	DeviceInfo string `json:"deviceInfo" db:"device_info"`
	IsActive   *bool  `json:"isActive"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}
