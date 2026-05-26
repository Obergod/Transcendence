package ws

type Hub struct {
	//registered client
	clients map[*Client]bool

	// msg from client
	broadcast chan []byte

	// register request from client
	register chan *Client

	// register request from client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		broadcast: make(chan []byte),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <-message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
