// I used the WebSocket example files for this, but I think this code is now mine? not sure.
// To feel a bit better I'm leaving this in the file:
// #region Copyright of examples used for this tiny inconsequential project
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// #endregion
package main

import (
	"WebApp/back"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")




func main() {
	flag.Parse() //I don't think this does anything?
	//region serve home
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
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
	})
	//endregion
	//region serve games
	http.HandleFunc("/{gamename}",func(w http.ResponseWriter, r *http.Request) {
		//validate path?
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "front/games/"+r.PathValue("gamename")+"/game.html")
	})
	//endregion
	//region serve game js
	http.HandleFunc("/js/{gamename}",func(w http.ResponseWriter, r *http.Request) {
		//validate path?
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFile(w, r, "front/games/"+r.PathValue("gamename")+"/game.js")
	})
	//endregion
	//region respond to WebSocket init requests
	http.HandleFunc("/ws/{gamename}", func(w http.ResponseWriter, r *http.Request) {
		back.ServeWs(w, r, r.PathValue("gamename"))
	})
	//endregion
	back.ManageGames() // start all game managers
	//run while you still can
	//err := http.ListenAndServeTLS(*addr, "weird/domain.crt", "weird/rootCA.key", nil)
	//never messing with certification again /hj. hell.
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
