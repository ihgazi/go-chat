package user

import (
    "time"
    "strconv"
    "context"

    "github.com/ihgazi/go-chat/util"
    "github.com/golang-jwt/jwt/v5"
    "github.com/ihgazi/go-chat/config"
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

// Custom JWT Claims class
type MyJWTClaims struct {
    ID string `json:"id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
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

    // Create JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
        ID: strconv.Itoa(int(u.ID)),
        Username: u.Username,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer: strconv.Itoa(int(u.ID)),
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24*time.Hour)), // expires in one day
        },
    })
    
    // Fetch secret key from environment
    // sign token with key
    secretKey := config.LoadEnv()
    ss, err := token.SignedString([]byte(secretKey))
    if err != nil {
        return &LoginUserRes{}, err
    }

    return &LoginUserRes{accessToken: ss, Username: u.Username, ID: strconv.Itoa(int(u.ID))}, nil
}

