package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AlexWilliam12/silent-signal/configs"
	"github.com/AlexWilliam12/silent-signal/handlers"
	"github.com/gorilla/mux"
)

var logger *configs.Logger

func init() {
	logger = configs.NewLogger("main")

	logger.Debug("Running initializers...")
	configs.Init()
	logger.Debug("Initalizers were finished successufully")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/login", handlers.HandleLogin).Methods("POST")
	r.HandleFunc("/auth/register", handlers.HandleRegister).Methods("POST")
	r.HandleFunc("/chat/private", handlers.HandlePrivateChat)
	r.HandleFunc("/chat/group", handlers.HandleGroupChat)

	port := ":" + os.Getenv("SERVER_PORT")

	logger.Infof("Server is running on port %s", port[1:])
	log.Fatal(http.ListenAndServe(port, r))
}
