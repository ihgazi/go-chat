package ws

import (
	"context"
	"database/sql"
	"strconv"
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
func (r *repository) JoinRoom(ctx context.Context, roomID string, userID string) error {
	query := `SELECT FROM room_member WHERE room_id = $1 AND user_id = $2`
	err := r.db.QueryRowContext(ctx, query, roomID, userID).Scan()

	if err == sql.ErrNoRows {
		query = `INSERT INTO room_member (room_id, user_id) VALUES ($1, $2)`
		_, err = r.db.ExecContext(ctx, query, roomID, userID)
	} else if err == nil {
		query = `UPDATE room_member SET last_online = NOW() WHERE room_id = $1 and user_id = $2`
		_, err = r.db.ExecContext(ctx, query, roomID, userID)
	} else {
		return err
	}

	return nil
}

// WriteMessage adds a new message to the room_message table
// It is called asynchronously with websocket messages
func (r *repository) WriteMessage(ctx context.Context, msg *Message) error {
	query := `INSERT INTO room_message (room_id, user_id, message) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, query, msg.RoomID, msg.UserID, msg.Content)
	if err != nil {
		return err
	}

	return nil
}

// FetchMessages retrieves messages for a specific room
// It is called when a user joins a room to load previous messages
func (r *repository) FetchRoomMessages(ctx context.Context, roomID string) ([]*Message, error) {
	query := `
        SELECT rm.user_id, u.username, rm.message
        FROM room_message rm
        JOIN users u ON rm.user_id = u.id
        WHERE rm.room_id = $1 AND
    	rm.created_at >= NOW() - INTERVAL '1 hour'
        ORDER BY rm.created_at ASC
    `
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.UserID, &msg.Username, &msg.Content); err != nil {
			return nil, err
		}
		msg.RoomID = roomID
		messages = append(messages, &msg)
	}
	return messages, nil
}

// Get the list of joined clients for a specific room.
func (r *repository) GetClients(ctx context.Context, roomID string) ([]ClientRes, error) {
	query := `
		SELECT u.id, u.username, u.is_online
		FROM users u
		JOIN room_member rm ON u.id = rm.user_id
		WHERE rm.room_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []ClientRes
	for rows.Next() {
		var c ClientRes
		if err := rows.Scan(&c.ID, &c.Username, &c.IsOnline); err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}
	return clients, nil
}

func (r *repository) SetUserOnlineStatus(ctx context.Context, userID string, isOnline bool) error {
	var query string
	if isOnline {
		query = `UPDATE users SET is_online = true, last_login = NOW() WHERE id = $1`
	} else {
		query = `UPDATE users SET is_online = false WHERE id = $1`
	}
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// Get entire list of rooms and their corresponding members. This is used to initialize the in-memory map maintained by Hub.
func (r *repository) GetRoomMembers(ctx context.Context) (map[string][]string, error) {
	query := `SELECT room_id, CAST(user_id AS VARCHAR) FROM room_member`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make(map[string][]string)
	for rows.Next() {
		var roomID, userID string
		if err := rows.Scan(&roomID, &userID); err != nil {
			return nil, err
		}
		members[roomID] = append(members[roomID], userID)
	}
	return members, nil
}
