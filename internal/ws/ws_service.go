package ws

import (
    "time"
    "context"
)

type service struct {
    Repository
    hub *Hub
    timeout time.Duration
}

func NewService(repository Repository, h *Hub) Service {
    return &service {
        repository,
        h,
        time.Duration(2) * time.Second,
    }
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

    room.Clients = make(map[string]*Client)
    s.hub.Rooms[room.ID] = room; 

    return &CreateRoomRes{
        ID: room.ID,
        Name: room.Name,
    }, nil
}

func (s *service) JoinRoom(c context.Context, cl *Client, m *Message) {
    // Register new client through the register channel
    s.hub.Register <- cl
    // Broadcast the message
    s.hub.Broadcast <- m

    go cl.WriteMessage()
    cl.ReadMessage(s.hub)  
}

func (s *service) GetRooms(ctx context.Context) (r []RoomRes) {
    var rooms []RoomRes
    for _, room := range s.hub.Rooms {
        rooms = append(rooms, RoomRes{
            ID: room.ID,
            Name: room.Name,
        })
    }

    return rooms;
}

func (s *service) GetClients(ctx context.Context, roomID string) (c []ClientRes) {
    var clients []ClientRes

	if _, ok := s.hub.Rooms[roomID]; !ok {
		clients = make([]ClientRes, 0)
	    return clients
	}

	for _, c := range s.hub.Rooms[roomID].Clients {
		clients = append(clients, ClientRes{
			ID: c.ID,
			Username: c.Username,
		})
	}

    return clients
}
