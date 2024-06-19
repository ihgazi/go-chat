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

// Login Request
type LoginUserReq struct {
    Email string `json:"email" db:"email"`
    Password string `json:"password" db:"password"`
}

// Login Response
type LoginUserRes struct {
    accessToken string
    Username string `json:"username" db:"username"`
    ID string `json:"id" db:"id"`
}

type Repository interface { 
    CreateUser(ctx context.Context, user *User) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type Service interface {
    CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error)
    Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error)
}
