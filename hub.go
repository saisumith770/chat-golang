package main

type ConnInfo struct{
	channelName string 
	client *Conn
}

type BroadCastMessage struct{
	channelName string 
	message []byte
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	channels map[string] []*Conn

	// Inbound messages from the clients.
	broadcast chan BroadCastMessage

	// Register requests from the clients.
	register chan ConnInfo

	// Unregister requests from clients.
	unregister chan ConnInfo
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan BroadCastMessage),
		register:   make(chan ConnInfo),
		unregister: make(chan ConnInfo),
		channels:   make(map[string] []*Conn),
	}
}

func remove(s []*Conn, i int) []*Conn {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func (h *Hub) run() {
	for {
		select {
		case req := <-h.register:
			h.channels[req.channelName] = append(h.channels[req.channelName],req.client)
		case req := <-h.unregister:
			for index,element := range h.channels[req.channelName]{
				if element == req.client{
					h.channels[req.channelName] = remove(h.channels[req.channelName],index)
					return 
				}
			}
			close(req.client.send)
		case broadcast := <-h.broadcast:
			for _,client := range h.channels[broadcast.channelName] {
				select {
				case client.send <- broadcast.message:
				default:
					h.unregister <- ConnInfo{
						channelName:broadcast.channelName,
						client: client,
					}
				}
			}
		}
	}
}
