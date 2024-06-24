package back

type Game interface { //idk...
	init() Game
	Gameloop()
	AddPlayer(*Client)
	RemovePlayer(*Client)
	playerInput([]byte, *Client)
}

// #region WordGame

type WordGame struct {
	players map[*Client]WordGamePlayer

	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func (b *WordGame) init() Game {
	instance := *new(WordGame)
	instance.players = make(map[*Client]WordGamePlayer)
	instance.broadcast = make(chan []byte)
	instance.register = make(chan *Client)
	instance.unregister = make(chan *Client)
	go instance.Gameloop()
	return &instance
}

type WordGamePlayer struct {
	client *Client
	ok     bool
}

func (b *WordGame) Gameloop() {
	for {
		select {
		case client := <-b.register:
			b.players[client] = WordGamePlayer{client: client, ok: true}
		case client := <-b.unregister:
			if _, ok := b.players[client]; ok {
				delete(b.players, client)
				close(client.send)
			}
		case message := <-b.broadcast:
			for client := range b.players {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(b.players, client)
				}
			}
		}
	}
}

func (b *WordGame) playerInput(msg []byte, c *Client) {
	b.broadcast <- msg
}
func (b *WordGame) AddPlayer(c *Client) {
	b.register <- c
}
func (b *WordGame) RemovePlayer(c *Client) {
	b.unregister <- c
}

//#endregion
