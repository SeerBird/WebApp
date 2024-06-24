// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package back

var managers map[string]AnyLobbyManager = map[string]AnyLobbyManager{
	"bababoi": LobbyManager[*WordGame]{games: make(map[*WordGame]bool)},
}

func ManageGames() {
	for manager := range managers {
		go managers[manager].start()
	}
}

func connectToGame(gamename string, args ...any) Game {
	manager := managers[gamename]
	if manager == nil {
		//scream
		return nil
	}
	return manager.connect(args)
}

type AnyLobbyManager interface {
	start()
	connect(args ...any) Game
}
type constraint interface {
	Game
	comparable
}
type LobbyManager[T constraint] struct {
	games    map[T]bool
	initiate chan T
	stop     chan T
}

func (m LobbyManager[T]) start() {
	for {
		select {
		case game := <-m.initiate:
			m.games[game] = true
		case game := <-m.stop:
			if _, ok := m.games[game]; ok {
				delete(m.games, game)
			}
		}
	}
}
func (m LobbyManager[T]) connect(args ...any) Game {
	game := (*new(T)).init()
	m.games[game.(T)] = true
	return game
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
