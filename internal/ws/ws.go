package ws

import (
	"context"
)

type CreateRoomReq struct {
	Name string `json:"name"`
}

type CreateRoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	IsOnline bool   `json:"is_online"`
}

type Repository interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	FetchRooms() ([]*Room, error)
	JoinRoom(ctx context.Context, roomID string, userID string) error
	WriteMessage(ctx context.Context, msg *Message) error
	FetchRoomMessages(ctx context.Context, roomID string) ([]*Message, error)
	GetClients(ctx context.Context, roomID string) ([]ClientRes, error)
	SetUserOnlineStatus(ctx context.Context, userID string, isOnline bool) error
	GetRoomMembers(ctx context.Context) (map[string][]string, error)
}

type Service interface {
	CreateRoom(ctx context.Context, req *CreateRoomReq) (*CreateRoomRes, error)
	Connect(ctx context.Context, cl *Client) error
	JoinRoom(ctx context.Context, roomID string, userID string) error
	GetRooms(ctx context.Context) (r []RoomRes)
	GetMyRooms(ctx context.Context, userID string) (r []RoomRes)
	GetClients(ctx context.Context, roomID string) (c []ClientRes)
	FetchRooms() error
}
