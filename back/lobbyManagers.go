package back
var managers map[string]AnyLobbyManager = map[string]AnyLobbyManager{
	"TicTacToe": &LobbyManager[*TTTGame, TTTGame]{games: make(map[*TTTGame]bool),
		stop: make(chan *TTTGame), join: make(chan JoinData)},
}

func ManageGames() {
	for manager := range managers {
		go managers[manager].start()
	}
}

func connectToGame(gamename string, client *Client, args string) {
	manager := managers[gamename]
	if manager == nil {
		//scream. handle pls.
		return
	}
	manager.connect(client, args)
}

type AnyLobbyManager interface {
	start()
	connect(client *Client, args string)
}

type ComparableGame[T any] interface {
	*T
	Game
	comparable
}
type LobbyManager[T ComparableGame[U], U any] struct {
	games map[T]bool
	stop  chan T
	join  chan JoinData
}

func (m *LobbyManager[T, U]) start() {
	for {
	out:
		select {
		case game := <-m.stop:
			delete(m.games, game)
		case data := <-m.join:
			for game := range m.games {
				if game.joinable(data.args) {
					game.addPlayer(data.client)
					break out
				}
			}
			//TODO: make sure game is joinable for data.args or handle otherwise
			game:=m.newGame()
			go game.gameloop()
			m.games[game] = true
			game.addPlayer(data.client)
		}
	}
}
func (m *LobbyManager[T, U]) newGame() T{
	var game U
	res:=(T(&game)).init(m)
	return (res).(T)
}

type JoinData struct {
	client *Client
	args   string
}

func (m *LobbyManager[T, U]) connect(client *Client, args string) {
	m.join <- JoinData{client: client, args: args}
}
