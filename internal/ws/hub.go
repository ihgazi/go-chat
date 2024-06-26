package ws

type Room struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Clients map[string]*Client `json:"clients"`
}

type Hub struct {
    Rooms map[string]*Room
    Register chan *Client
    Unregister chan *Client
    Broadcast chan *Message
}

func NewHub() *Hub {
    return &Hub{
        Rooms: make(map[string]*Room),
        Register: make(chan *Client),
        Unregister: make(chan *Client),
        Broadcast: make(chan *Message, 5),

    }
}

func (h *Hub) Run() {
    for {
        select {
        case cl := <-h.Register:
            if _, ok := h.Rooms[cl.RoomID]; ok {
                r := h.Rooms[cl.RoomID]

                if _, ok := r.Clients[cl.ID]; !ok {
                    r.Clients[cl.ID] = cl
                }
            }
        case cl := <-h.Unregister:
            if _, ok := h.Rooms[cl.RoomID]; ok {
                r := h.Rooms[cl.RoomID]

                if _, ok := r.Clients[cl.ID]; ok {
                    // Broadcast exit message
                    h.Broadcast <- &Message{
                        Content: "user has left the chat",
                        RoomID: cl.RoomID,
                        Username: cl.Username,
                        UserID: cl.ID,
                    }

                    delete(r.Clients, cl.ID)
                    close(cl.Message)
                }
            }
        case msg := <-h.Broadcast:
            if _, ok := h.Rooms[msg.RoomID]; ok {
                for _, cl := range h.Rooms[msg.RoomID].Clients {
                    cl.Message <- msg
                }
            }
        }
    }   
}
