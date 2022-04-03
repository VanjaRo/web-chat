package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/VanjaRo/web-chat/config"
	"github.com/VanjaRo/web-chat/repositories"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	flag.Parse()

	config.InitRedis()
	db := config.InitDB()

	wsServer := NewWsServer(&repositories.RoomRepository{Db: db}, &repositories.UserRepository{Db: db})
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
