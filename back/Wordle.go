package back

import (
	"log"
	"math"
	"math/rand"
)

const size int = 5

var alphabetValues map[string]int = map[string]int{
	"A": 1,
	"B": 1,
	"C": 1,
	"D": 1,
	"E": 1,
	"F": 1,
	"G": 1,
	"H": 1,
	"I": 1,
	"J": 1,
	"K": 1,
	"L": 1,
	"M": 1,
	"N": 1,
	"O": 1,
	"P": 1,
	"Q": 1,
	"R": 1,
	"S": 1,
	"T": 1,
	"U": 1,
	"V": 1,
	"W": 1,
	"X": 1,
	"Y": 1,
	"Z": 1,
}

type Wordle struct {
	players map[*Client]*WordlePlayer

	grid [size][size]string
	turn int

	inPacketChannel chan *ClientPacket
	unregister      chan *Client
	register        chan *Client
	manager         AnyLobbyManager
}
type WordlePlayer struct {
	client *Client
	score  int
	order  int //0,1 are team1 and 2,3 are team2
}

func (b *Wordle) init(m AnyLobbyManager) Game {
	inst := *b
	inst.players = make(map[*Client]*WordlePlayer)
	inst.inPacketChannel = make(chan *ClientPacket)
	inst.unregister = make(chan *Client)
	inst.register = make(chan *Client)
	inst.manager = m
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			inst.grid[i][j] = ""
		}
	}
	return &inst
}

// region gameloop
func (b *Wordle) gameloop() {
	for {
		select {
		case client := <-b.register:
			b.players[client] = &WordlePlayer{client: client, score: 0, order: len(b.players)}
			go client.readPump(b)
			go client.writePump()
			if len(b.players) == 4 {
				b.grid = newGrid()
				b.turn = 0
				for anyClient := range b.players {
					anyClient.send(b.getServerPacket(b.players[anyClient]), "update")
				}
			} else {
				client.send(b.getServerPacket(b.players[client]), "update")
			}
		case client := <-b.unregister:
			if _, ok := b.players[client]; ok {
				delete(b.players, client)
				close(client.sendChannel)
				for client := range b.players {
					client.send(b.getServerPacket(b.players[client]), "update")
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
func newGrid() [size][size]string {
	var grid [size][size]string
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			var alphabet []string
			for letter := range alphabetValues {
				alphabet = append(alphabet, letter)
			}
			grid[i][j] = alphabet[int(math.Floor(rand.Float64()*float64(len(alphabet))))]
		}
	}
	return grid
}
func (b *Wordle) handleInput(packet *ClientPacket) { //this is where the magic happens
	player := b.players[packet.client]
	if player.order != b.turn {
		return
	}
	if packet.msg["type"] != "input" {
		log.Output(0, "Unknown ClientPacket type: "+(packet.msg["type"]).(string))
		return
	}
	msg := packet.msg["value"].([]map[string]int)
	//region get the word the player made
	var letterCoords []coord
	for _, letter := range msg {
		thisCoord := coord{i: letter["i"], j: letter["j"]}
		letterCoords = append(letterCoords, thisCoord)
		//region validate
		if thisCoord.i < 0 || thisCoord.i > size-1 || thisCoord.j < 0 || thisCoord.j > size-1 {
			log.Output(0, "Invalid letter coord input")
			return
		}
		if len(letterCoords) != 1 {
			prevCoord := letterCoords[len(letterCoords)-1]
			if math.Abs(float64(thisCoord.i-prevCoord.i)) > 1 || math.Abs(float64(thisCoord.j-prevCoord.j)) > 1 {
				log.Output(0, "Discontinuous letter coord input sequence")
				return
			}
		}
		//endregion
	}
	word := ""
	for _, coord := range letterCoords {
		word += b.grid[coord.i][coord.j]
	}
	//endregion
	
}

type coord struct {
	i int
	j int
}

// endregion
// region ServerPacket
func (b *Wordle) getServerPacket(recipient *WordlePlayer) ServerPacket { //grid, turn, playerList, index of player this is sent to(determined later)
	var packet ServerPacket = ServerPacket{
		grid:        b.grid,
		turn:        b.turn,
		playerList:  make(map[int]playerData),
		clientOrder: recipient.order,
	}
	for client := range b.players {
		player := b.players[client]
		packet.playerList[player.order] = playerData{name: client.conn.LocalAddr().String(), score: player.score}
	}
	return packet
}

type ServerPacket struct {
	grid        [size][size]string
	turn        int
	clientOrder int
	playerList  map[int]playerData
}
type playerData struct {
	name  string
	score int
}

// endregion
// region external stuff
func (b *Wordle) receivePacket(message map[string]interface{}, c *Client) {
	b.inPacketChannel <- &ClientPacket{msg: message, client: c}
}
func (b *Wordle) addPlayer(client *Client) {
	b.register <- client
}
func (b *Wordle) removePlayer(client *Client) {
	b.unregister <- client
}
func (b *Wordle) joinable(args string) bool {
	return len(b.players) < 4
}

//endregion
