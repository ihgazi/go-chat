package ws

import (
    "context"
    "strconv"
    "database/sql"
)

type DBTX interface {
    ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
    PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
    QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
    Query(query string, args ...interface{}) (*sql.Rows, error)
}

type repository struct {
    db DBTX
}

func NewRepository(db DBTX) Repository {
    return &repository{db: db}   
}

func (r *repository) CreateRoom(ctx context.Context, room *Room) (*Room, error) {
    var lastInsertID int
    query := `INSERT INTO room (name) VALUES ($1) RETURNING id`
    err := r.db.QueryRowContext(ctx, query, room.Name).Scan(&lastInsertID)

    if err != nil {
        return &Room{}, err
    }

    room.ID = strconv.Itoa(lastInsertID)
    return room, nil    
}

func (r *repository) FetchRooms() ([]*Room, error) {
    query := `SELECT id, name FROM room`
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var rooms []*Room
    for rows.Next() {
        var room Room
        if err := rows.Scan(&room.ID, &room.Name); err != nil {
            return nil, err
        }
        rooms = append(rooms, &room)
    }

    return rooms, nil
}
