package user

import (
    "time"
    "strconv"
    "context"

    "github.com/ihgazi/go-chat/util"
)


// User service takes CreateUserReq object from API endpoint
// Creates a User struct and passes it to repository

type service struct {
    Repository
    timeout time.Duration
}

func NewService(repository Repository) Service {
    return &service{
        repository,
        time.Duration(2) * time.Second,
    }
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
    ctx, cancel := context.WithTimeout(c, s.timeout)
    defer cancel()

    hashedPassword, err := util.HashPassword(req.Password)
    u := &User{
        Username: req.Username,
        Email: req.Email,
        Password: hashedPassword,
    }

    r, err := s.Repository.CreateUser(ctx, u)
    if err != nil {
        return nil, err
    }

    res := &CreateUserRes{
        ID: strconv.Itoa(int(r.ID)),
        Username: r.Username,
        Email: r.Email,
    }

    return res, nil;
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
    ctx, cancel := context.WithTimeout(c, s.timeout)
    defer cancel()

    // Call repostory method to get user details
    u, err := s.Repository.GetUserByEmail(ctx, req.Email) 
    if err != nil {
        return &LoginUserRes{}, err
    }

    // Check if password matches
    err = util.CheckPassword(req.Password, u.Password)
    if err != nil {
        return &LoginUserRes{}, err
    }

    ss, err := util.JWTCreateToken(strconv.Itoa(int(u.ID)), u.Username)

    if err != nil {
        return &LoginUserRes{}, err
    }

    return &LoginUserRes{accessToken: ss, Username: u.Username, ID: strconv.Itoa(int(u.ID))}, nil
}

