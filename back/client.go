package back

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the game.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	sendChannel chan ServerMessage
}
type ServerMessage struct{
	Msg any
	Tag string
}
// readPump pumps messages from the websocket connection to the game.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
// There is a bunch of code that I don't think I care about here, but it's all good-practise handling so if I use this later
// I should at least read the error handling code from the example. Therefore, it stays.
func (c *Client) readPump(game Game) {
	defer func() {
		game.removePlayer(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, packet, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		packet = bytes.TrimSpace(bytes.Replace(packet, newline, space, -1))
		var message map[string]interface{}
		err = json.Unmarshal(packet, &message)
		if err != nil {
			panic(err)
		}
		game.receivePacket(message, c)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.sendChannel:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			packet, err := json.Marshal(message)
			//handle err?
			w.Write(packet)

			// Add queued chat messages to the current websocket message.
			n := len(c.sendChannel)
			for i := 0; i < n; i++ {
				w, err := c.conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				packet, err = json.Marshal(<-c.sendChannel)
				//handle err?
				w.Write(packet)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (c *Client) send(msg any, tag string) {
	c.sendChannel<-ServerMessage{Msg:msg,Tag:tag}
}

// serveWs handles websocket requests from the peer.
func ServeWs(w http.ResponseWriter, r *http.Request, gamename string) {
	// use parameters in the request to determine what kind of game this starts/connects the user to
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	connectToGame(gamename, &Client{conn: conn, sendChannel: make(chan ServerMessage)}, "")
}
