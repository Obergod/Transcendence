package ws

type Hub struct {
	// L'annuaire : associe l'ID du joueur à sa connexion WebSocket
	clients map[uint]*Client

	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Si le joueur est déjà connecté ailleurs (ex: un autre onglet), on ferme l'ancien
			if oldClient, ok := h.clients[client.UserID]; ok {
				close(oldClient.send)
			}
			h.clients[client.UserID] = client

		case client := <-h.unregister:
			// On s'assure qu'on supprime bien le bon client
			if currentClient, ok := h.clients[client.UserID]; ok && currentClient == client {
				delete(h.clients, client.UserID)
				close(client.send)
			}

		case message := <-h.broadcast:
			// Pour l'instant on garde le broadcast, mais on pourra facilement cibler un ID précis plus tard !
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.UserID)
				}
			}
		}
	}
}