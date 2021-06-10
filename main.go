package main

import (
	"chat-app/config"
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	mongoConn = config.MongoConnection()
	flag.Parse()
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
