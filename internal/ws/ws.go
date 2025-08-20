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
}

type Repository interface {
	CreateRoom(ctx context.Context, room *Room) (*Room, error)
	FetchRooms() ([]*Room, error)
	JoinRoom(ctx context.Context, client *Client) error
	WriteMessage(ctx context.Context, msg *Message) error
}

type Service interface {
	CreateRoom(ctx context.Context, req *CreateRoomReq) (*CreateRoomRes, error)
	JoinRoom(ctx context.Context, cl *Client, m *Message) error
	GetRooms(ctx context.Context) (r []RoomRes)
	GetClients(ctx context.Context, roomID string) (c []ClientRes)
	FetchRooms() error
}
