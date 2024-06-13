package user

import (
    "time"
    "strconv"
    "context"
)

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

    // TODO: Implement password hashing
    u := &User{
        Username: req.Username,
        Email: req.Email,
        Password: req.Password,
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
