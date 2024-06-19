package user

import (
    "context"
)

// User model
type User struct {
    ID int64 `json:"id" db:"id"`
    Username string `json:"username" db:"username"`
    Email string `json:"email" db:"email"`
    Password string `json:"password" db:"password"`
}

// Create User Request 
type CreateUserReq struct {
    Username string `json:"username" db:"username"`
    Email string `json:"email" db:"email"`
    Password string `json:"password" db:"password"`
}

// Create User Response
type CreateUserRes struct {
    ID string `json:"id" db:"id"`
    Username string `json:"username" db:"username"`
    Email string `json:"email" db:"email"`
}

type Repository interface { 
    CreateUser(ctx context.Context, user *User) (*User, error)
}

type Service interface {
    CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
}
