package users

import "time"

type Role string



const (
    RoleMember  Role = "Member"
    RoleAdmin   Role = "Admin"
    RoleModerator Role = "Moderator"
    RoleGuest   Role = "Guest"
)

type UserData struct {
	UserID    string   		`json:"userId"`
    Username string  		`json:"username"`
    Email    string  		`json:"email"`
    Password *string 		`json:"password,omitempty"`
    Role     Role    		`json:"role"`
	IsDeleted bool        	`json:"isDeleted"`
	Balance  int    		`json:"balance"`
	CreatedAt  time.Time  	`json:"created_at"`
	UpdatedAt  *time.Time  	`json:"updated_at,omitempty"`
}


type CreateUSer struct {
	Username  string  `json:"username"`
	Email    string  `json:"email"`
	Password   *string  `json:"password"`
	Role   Role    `json:"role"`
	Balance   int    	`json:"balance"`
}

type UpdateUser   struct {
	UserID    string  `json:"userId"`
	Username  string  `json:"username"`
	Email    string  `json:"email"`
	Role   Role    `json:"role"`
}


type CreateDepositUser struct {
	UserID   string   `json:"userId"`
	Username   string  `json:"username"`
	PaymentCode  string   `json:"paymentCode"`
	Amount     int   `json:"amount"`
}

type DeleteUser struct  {
	Username   string  `json:"username"`
	IsDeleted   bool   `json:"isDeleted"`
}