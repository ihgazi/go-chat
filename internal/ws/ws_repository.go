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

// JoinRoom adds a new entry to room_member table
// if user already exists update last_online time
func (r *repository) JoinRoom(ctx context.Context, client *Client) error {
    query := `SELECT FROM room_member WHERE room_id = $1 AND user_id = $2`
    err := r.db.QueryRowContext(ctx, query, client.RoomID, client.ID).Scan()

    if err == sql.ErrNoRows {
        query = `INSERT INTO room_member (room_id, user_id) VALUES ($1, $2)`
        _, err = r.db.ExecContext(ctx, query, client.RoomID, client.ID)
    } else if err == nil {
        query = `UPDATE room_member SET last_online = NOW() WHERE room_id = $1 and user_id = $2`
        _, err = r.db.ExecContext(ctx, query, client.RoomID, client.ID)
    } else {
        return err;
    }

    return nil;
}
