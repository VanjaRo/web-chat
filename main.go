package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/VanjaRo/web-chat/auth"
	"github.com/VanjaRo/web-chat/config"
	"github.com/VanjaRo/web-chat/repositories"
)

var addr = flag.String("addr", ":8080", "http server address")

func main() {
	flag.Parse()

	config.InitRedis()
	db := config.InitDB()

	userRepository := &repositories.UserRepository{Db: db}

	wsServer := NewWsServer(&repositories.RoomRepository{Db: db}, userRepository)
	go wsServer.Run()

	api := &Api{
		UserRepository: userRepository,
	}

	http.HandleFunc("/ws", auth.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	}))
	http.HandleFunc("/api/login", api.LoginHandler)

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
