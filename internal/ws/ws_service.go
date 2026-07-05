package ws

import (
	"context"
	"log"
	"time"
)

type service struct {
	Repository
	hub     *Hub
	timeout time.Duration
}

func NewService(repository Repository, h *Hub) Service {
	s := &service{
		repository,
		h,
		time.Duration(2) * time.Second,
	}

	err := s.FetchRooms()
	if err != nil {
		panic(err)
	}

	return s
}

// Populate hub with rooms from database
func (s *service) FetchRooms() error {
	rooms, err := s.Repository.FetchRooms()
	if err != nil {
		return err
	}

	membersMap, err := s.Repository.GetRoomMembers(context.Background())
	if err != nil {
		// Log the error but don't crash, it might be an empty db initially
		log.Printf("Error fetching room members: %v", err)
	}

	for _, room := range rooms {
		room.Members = make(map[string]bool)
		if members, ok := membersMap[room.ID]; ok {
			for _, userID := range members {
				room.Members[userID] = true
			}
		}
		s.hub.Rooms[room.ID] = room
	}

	return nil
}

func (s *service) CreateRoom(c context.Context, req *CreateRoomReq) (*CreateRoomRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	r := &Room{
		Name: req.Name,
	}

	room, err := s.Repository.CreateRoom(ctx, r)
	if err != nil {
		return nil, err
	}

	room.Members = make(map[string]bool)
	s.hub.Rooms[room.ID] = room

	return &CreateRoomRes{
		ID:   room.ID,
		Name: room.Name,
	}, nil
}

func (s *service) Connect(c context.Context, cl *Client) error {
	// Register new client through the register channel
	s.hub.Register <- cl

	go cl.WriteMessage()
	cl.ReadMessage(s.hub)

	return nil
}

func (s *service) JoinRoom(c context.Context, roomID string, userID string) error {
	err := s.Repository.JoinRoom(c, roomID, userID)
	if err != nil {
		return err
	}

	// Update the in-memory cache
	if room, ok := s.hub.Rooms[roomID]; ok {
		room.Members[userID] = true
	}

	return nil
}

func (s *service) GetRooms(ctx context.Context) (r []RoomRes) {
	var rooms []RoomRes
	for _, room := range s.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   room.ID,
			Name: room.Name,
		})
	}

	return rooms
}

func (s *service) GetMyRooms(ctx context.Context, userID string) (r []RoomRes) {
	var myRooms []RoomRes

	for _, room := range s.hub.Rooms {
		if _, ok := room.Members[userID]; ok {
			myRooms = append(myRooms, RoomRes{
				ID:   room.ID,
				Name: room.Name,
			})
		}
	}

	if myRooms == nil {
		myRooms = make([]RoomRes, 0)
	}

	return myRooms
}

func (s *service) GetClients(ctx context.Context, roomID string) (c []ClientRes) {
	clients, err := s.Repository.GetClients(ctx, roomID)
	if err != nil || clients == nil {
		return make([]ClientRes, 0)
	}

	return clients
}
