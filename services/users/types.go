package users

import "time"

type Role string

const (
	RoleMember    Role = "Member"
	RoleAdmin     Role = "Admin"
	RoleModerator Role = "Moderator"
	RoleGuest     Role = "Guest"
)

type UserData struct {
	UserID    string     `json:"userId" db:"user_id"`
	FullName  *string    `json:"fullName"  db:"full_name"`
	Username  string     `json:"username"  db:"username"`
	Email     string     `json:"email"  db:"email"`
	Password  *string    `json:"password,omitempty"  db:"password"`
	Role      Role       `json:"role"  db:"role"`
	IsDeleted bool       `json:"isDeleted"  db:"is_deleted"`
	Balance   int        `json:"balance" db:"balance"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type CreateUser struct {
	FullName *string `json:"fullName"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password *string `json:"password"`
	Role     Role    `json:"role"`
	Balance  int     `json:"balance"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUser struct {
	FullName *string `json:"fullName"`
	UserID   string  `json:"userId"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Role     Role    `json:"role"`
}

type CreateDepositUser struct {
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	PaymentCode string `json:"paymentCode"`
	Amount      int    `json:"amount"`
}

type DeleteUser struct {
	Username  string `json:"username"`
	IsDeleted bool   `json:"isDeleted"`
}
