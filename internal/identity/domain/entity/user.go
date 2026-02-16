package entity

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
	ID        string     `db:"id"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	Email     string     `db:"email"`
	Role      UserRole   `db:"role"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func NewUser(params NewUserParams) *User {
	return &User{
		ID:        params.ID,
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Email:     params.Email,
		Role:      UserRoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
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
