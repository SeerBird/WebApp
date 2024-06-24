package back

type Game interface { //idk...
	init()
	Gameloop()
	AddPlayer()
	RemovePlayer(Client)
}

// #region WordGame

type WordGame struct {
	players map[WordGamePlayer]*Client
}

type WordGamePlayer struct {
}

func (b WordGame) init() {

}

func (b WordGame) Gameloop() {

}
func (b WordGame) AddPlayer() {

}
func (b WordGame) RemovePlayer(Client) {

}

//#endregion
