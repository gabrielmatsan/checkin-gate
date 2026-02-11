package identity

import "time"

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

type NewUserParams struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
}

type User struct {
	ID        string     `json:"id" db:"id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Email     string     `json:"email" db:"email"`
	Role      UserRole   `json:"role" db:"role" default:"user"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

func NewUser(params NewUserParams) *User {
	return &User{
		ID:        params.ID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Role:      UserRoleUser,
		CreatedAt: time.Now(),
	}
}

func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

func (u *User) PromoteToAdmin() {
	u.Role = UserRoleAdmin
	u.touch()
}

func (u *User) touch() {
	now := time.Now()
	u.UpdatedAt = &now
}

