package main

import (
	"chat-app/config"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	var addr = flag.String("addr", fmt.Sprintf(":%v", os.Getenv("PORT")), "http server address")
	mongoConn = config.MongoConnection()
	flag.Parse()
	wsServer := NewWebsocketServer()
	go wsServer.Run()

	http.HandleFunc("/websocket", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})
	http.HandleFunc("/upload", uploadFile)

	http.Handle("/images", http.FileServer(http.Dir("./public/images")))

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
