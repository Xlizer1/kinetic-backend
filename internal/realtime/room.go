package realtime

type Room struct {
	ID        uint
	Type      string
	Clients   map[*Client]bool
	Broadcast chan []byte
}

func NewRoom(id uint, roomType string) *Room {
	return &Room{
		ID:        id,
		Type:      roomType,
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan []byte, 256),
	}
}

func (r *Room) Join(client *Client) {
	r.Clients[client] = true
	client.Rooms[r.ID] = true
}

func (r *Room) Leave(client *Client) {
	if _, ok := r.Clients[client]; ok {
		delete(r.Clients, client)
		delete(client.Rooms, r.ID)
	}
}

func (r *Room) BroadcastMessage(message []byte) {
	for client := range r.Clients {
		select {
		case client.Send <- message:
		default:
			r.Leave(client)
		}
	}
}

func (r *Room) Count() int {
	return len(r.Clients)
}
