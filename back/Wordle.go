package back

import "fmt"

import (
	"log"
	"math"
	"math/rand"
	"os/exec"
	"slices"
	"strings"
	"os"
	"path/filepath"
)

const size int = 5
const playercount int = 4

var letterFreq map[string]float64 = map[string]float64{
	"A": 8.2,
	"B": 1.5,
	"C": 2.8,
	"D": 4.3,
	"E": 12.7,
	"F": 2.2,
	"G": 2.0,
	"H": 6.1,
	"I": 7.0,
	"J": 0.15,
	"K": 0.77,
	"L": 4.0,
	"M": 2.4,
	"N": 6.7,
	"O": 7.5,
	"P": 1.9,
	"Q": 0.095,
	"R": 6.0,
	"S": 6.3,
	"T": 9.1,
	"U": 2.8,
	"V": 0.98,
	"W": 2.4,
	"X": 0.15,
	"Y": 2.0,
	"Z": 0.074,
}
var alphabetValues map[string]int = map[string]int{
	"A": 1,
	"B": 4,
	"C": 5,
	"D": 3,
	"E": 1,
	"F": 5,
	"G": 3,
	"H": 4,
	"I": 1,
	"J": 7,
	"K": 6,
	"L": 3,
	"M": 4,
	"N": 2,
	"O": 1,
	"P": 4,
	"Q": 8,
	"R": 2,
	"S": 2,
	"T": 2,
	"U": 4,
	"V": 5,
	"W": 5,
	"X": 7,
	"Y": 4,
	"Z": 8,
}

type Wordle struct {
	players map[*Client]*WordlePlayer

	grid    [size][size]string
	turn    int
	doubleW coord
	doubleL coord

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
	inst.turn = -1
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			inst.grid[i][j] = "?"
		}
	}
	return &inst
}

// region gameloop
func (b *Wordle) gameloop() {
	for {
		select {
		case client := <-b.register:
			flag := true
			i := 0
			for ;flag;i++ {
				flag=false
				for _, player := range b.players {
					if i == player.order {
						flag=true
					}
				}
			}
			i--
			if i>playercount-1{
				panic("player reg went wrong")
			}
			b.players[client] = &WordlePlayer{client: client, score: 0, order: i}
			go client.readPump(b)
			go client.writePump()
			if len(b.players) == playercount {
				//region start game
				b.regenerate()
				b.turn = 0
				//region swap teams half the time
				if rand.Float64() > 0.5 {
					for _, player := range b.players {
						switch player.order {
						case 0:
							player.order = 2
						case 1:
							player.order = 3
						case 2:
							player.order = 0
						case 3:
							player.order = 1
						}
					}
				}
				//endregion
				for anyClient := range b.players {
					anyClient.send(b.getServerPacket(b.players[anyClient]), "update")
				}
				//endregion
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
	var sum float64 = 0
	var alphabet []string
	for letter := range alphabetValues {
		sum += letterFreq[letter]
		alphabet = append(alphabet, letter)
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			rand := rand.Float64() * sum
			grid[i][j] = "?"
			for _, letter := range alphabet {
				rand -= letterFreq[letter]
				if rand <= 0 {
					grid[i][j] = letter
					break
				}
			}
		}
	}
	return grid
}
func (b *Wordle) handleInput(packet *ClientPacket) { //this is where the magic happens
	//region get and validate the player
	player := b.players[packet.client]
	if player.order != b.turn {
		return
	}
	if packet.msg["type"] != "input" {
		log.Output(0, "Unknown ClientPacket type: "+(packet.msg["type"]).(string))
		return
	}
	//endregion
	//region get the word the player made
	var letterCoords []coord
	for _, uncastLetter := range packet.msg["value"].([]any) {
		letter := uncastLetter.(map[string]any)
		thisCoord := coord{i: int(letter["i"].(float64)), j: int(letter["j"].(float64))}
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
	if checkWord(word){
		//region give the points
		var points int = 0
		for _, letter := range letterCoords {
			add := alphabetValues[b.grid[letter.i][letter.j]]
			if letter == b.doubleL {
				add *= 2
			}
			points += add
		}
		if slices.Contains(letterCoords, b.doubleW) {
			points *= 2
		}
		player.score += points
		//endregion
		//region replace the letters
		newGrid := newGrid()
		for _, coord := range letterCoords {
			b.grid[coord.i][coord.j] = newGrid[coord.i][coord.j]
		}
		//endregion
		//region increment turn and restart the round if it's over
		b.turn = (b.turn + 1) % playercount
		if b.turn == 0 {
			b.regenerate()
		}
		//endregion
		for client := range b.players {
			client.send(b.getServerPacket(b.players[client]), "update")
		}
	}
}
func (b *Wordle) regenerate() {
	b.grid = newGrid()
	b.doubleL = coord{i: randCoord(), j: randCoord()}
	b.doubleW = coord{i: randCoord(), j: randCoord()}
}
func checkWord(word string) bool {
	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("python", filepath.Dir(path)+"\\back\\resources\\Wordle\\checkWord.py", word)
	var out strings.Builder
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Output(0, err.Error())
	}
	output:=out.String()
	log.Output(0,output)
	return out.String() == "True\r\n"
}
func randCoord() int {
	return int(rand.Float64() * float64(size))
}

type coord struct {
	i int
	j int
}

// endregion
// region ServerPacket
func (b *Wordle) getServerPacket(recipient *WordlePlayer) ServerPacket { //grid, turn, playerList, index of player this is sent to(determined later)
	var packet ServerPacket = ServerPacket{
		Grid:        b.grid,
		Turn:        b.turn,
		ClientOrder: recipient.order,
	}
	for client := range b.players {
		player := b.players[client]
		packet.PlayerList[player.order] = playerData{Name: fmt.Sprint(player.order), Score: player.score}
	}
	return packet
}

type ServerPacket struct {
	Grid        [size][size]string
	Turn        int
	ClientOrder int
	PlayerList  [playercount]playerData
}
type playerData struct {
	Name  string
	Score int
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
