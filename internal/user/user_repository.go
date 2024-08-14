package user

import (
    "errors"
    "context"
    "database/sql"
)

// User repository injected with database connection object
// Takes a User struct and updates the database

type DBTX interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type repository struct {
    db DBTX
}

func NewRepository(db DBTX) Repository {
    return &repository{db: db}   
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
    // Check if user is already registered
    query := `SELECT id FROM users WHERE email = $1`
    rows, err := r.db.QueryContext(ctx, query, user.Email)
    if rows.Next() {
        return &User{}, errors.New("User already exists!")
    }

    var lastInsertID int
    query = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
    err = r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(&lastInsertID)
    
    if err != nil {
        return &User{}, err
    }

    user.ID = int64(lastInsertID)
    return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    u := User{}
    query := `SELECT id, username, email, password FROM users WHERE email = $1`
    err := r.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
    
    if err != nil {
        return &User{}, err
    }

    return &u, nil
}

