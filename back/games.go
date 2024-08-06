package back

import (
	"log"
)

type Game interface {
	init(AnyLobbyManager) Game
	gameloop()
	addPlayer(*Client) //never called concurrently
	removePlayer(*Client)
	receivePacket(map[string]interface{}, *Client)
	joinable(string) bool
}
type ClientPacket struct {
	msg    map[string]interface{}
	client *Client
}

// #region TTTGame

type TTTGame struct {
	players map[*Client]*TTTGamePlayer

	inPacketChannel chan *ClientPacket

	grid        *[3][3]int
	currentRole int //1 for cross, 0 for circle
	// Unregister requests from clients.
	unregister chan *Client
	register   chan *Client
	manager    AnyLobbyManager
}
type TTTGamePlayer struct {
	client *Client
	score  int
	role   int //1 for cross, 0 for circle
}

func (b *TTTGame) init(m AnyLobbyManager) Game {
	instance := *b
	instance.players = make(map[*Client]*TTTGamePlayer)
	instance.inPacketChannel = make(chan *ClientPacket)
	instance.unregister = make(chan *Client)
	instance.register = make(chan *Client)
	instance.manager = m
	instance.resetGrid()
	return &instance
}

// #region gameplay
func (b *TTTGame) handleInput(packet *ClientPacket) {
	msg := packet.msg["value"].(map[string]interface{})
	if b.currentRole != b.players[packet.client].role {
		//wrong turn
		return
	}
	i := int(msg["i"].(float64))
	j := int(msg["j"].(float64))
	if b.grid[i][j] != -1 {
		//bad input
		return
	}
	b.grid[i][j] = b.players[packet.client].role
	for client, _ := range b.players {
		client.send(*b.grid, "update")
	}
	if b.checkWin(b.currentRole) {
		b.resetGrid()
		for client, _ := range b.players {
			client.send(*b.grid, "update")
		}
	}
	if b.currentRole == 0 {
		b.currentRole = 1
	} else {
		b.currentRole = 0
	}
}
func (b *TTTGame) resetGrid() {
	b.grid = &[3][3]int{{-1, -1, -1}, {-1, -1, -1}, {-1, -1, -1}}
}
func (b *TTTGame) checkWin(role int) bool {
	win := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			win = win && (b.grid[i][j] == role)
		}
		if win {
			return true
		}
		win = true
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			win = win && (b.grid[j][i] == role)
		}
		if win {
			return true
		}
		win = true
	}

	for i := 0; i < 3; i++ {
		win = (b.grid[i][i] == role) && win
	}
	if win {
		return true
	}
	win = true
	for i := 0; i < 3; i++ {
		win = (b.grid[i][2-i] == role) && win
	}
	if win {
		return true
	}
	end := true
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			end = end && (b.grid[i][j] != -1)
		}
	}
	return end
}

// #endregion
func (b *TTTGame) gameloop() {
	for {
		select {
		case client := <-b.register:
			b.players[client] = &TTTGamePlayer{client: client, role: len(b.players), score: 0}
			go client.readPump(b)
			go client.writePump()
			client.send(*b.grid, "update")
		case client := <-b.unregister:
			if _, ok := b.players[client]; ok {
				delete(b.players, client)
				close(client.sendChannel)
				b.resetGrid()
				for client, _ := range b.players {
					client.send(*b.grid, "update")
				}
			}
		case packet := <-b.inPacketChannel:
			msg := packet.msg
			if msgType, ok := msg["type"]; ok {
				switch msgType {
				case "input":
					b.handleInput(packet)
				default:
					log.Output(0, "Unexpected message type received")
				}
			} else {
				log.Output(0, "Wrong message format received")
			}
		}
	}
}

func (b *TTTGame) receivePacket(message map[string]interface{}, c *Client) {
	b.inPacketChannel <- &ClientPacket{msg: message, client: c}
}
func (b *TTTGame) addPlayer(client *Client) {
	b.register <- client
}
func (b *TTTGame) removePlayer(client *Client) {
	b.unregister <- client
}
func (b *TTTGame) joinable(args string) bool {
	return len(b.players) < 2
}

//#endregion
