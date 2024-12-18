package back

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

