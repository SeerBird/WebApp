module WebApp

go 1.22.4

require github.com/gorilla/websocket v1.5.2

require golang.org/x/net v0.25.0 // indirect
//there was a bunch of indirect dependencies here but 'go mod tidy' somehow got rid of them even though deleting them broke everything?
//I don't get it.