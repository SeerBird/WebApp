// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"WebApp/back"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "front/home.html")
}

func main() {
	flag.Parse() //I don't think this does anything?
	//give access to front files
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("front"))))
	http.HandleFunc("/", serveHome)
	//define multiplayer connection behaviour
	http.HandleFunc("/ws/{gamename}", func(w http.ResponseWriter, r *http.Request) {
		back.ServeWs(w, r, r.PathValue("gamename"))
	})
	//manage all the game lobbies
	back.ManageGames()
	//run while you still can
	//err := http.ListenAndServeTLS(*addr, "weird/domain.crt", "weird/rootCA.key", nil)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
