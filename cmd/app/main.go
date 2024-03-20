package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AlexWilliam12/silent-signal/internal/configs"
	"github.com/AlexWilliam12/silent-signal/internal/handlers"
	"github.com/gorilla/mux"
)

func init() {
	log.Println("Running initializers...")
	configs.Init()
	log.Println("Initalizers were finished successufully")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/auth/login", handlers.HandleLogin).Methods("POST")
	r.HandleFunc("/auth/register", handlers.HandleRegister).Methods("POST")

	r.HandleFunc("/user", handlers.HandleUserUpdate).Methods("PUT")
	r.HandleFunc("/user", handlers.HandleUserDelete).Methods("DELETE")
	r.HandleFunc("/user/contact", handlers.HandleSaveContact).Methods("POST")

	r.HandleFunc("/upload/picture", handlers.HandleFetchPicture).Methods("GET")
	r.HandleFunc("/upload/picture", handlers.HandleUploadPicture).Methods("POST", "PUT")
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads")))).Methods("GET")

	r.HandleFunc("/chat/private", handlers.HandlePrivateConnections)
	go handlers.HandlePrivateMessages()

	r.HandleFunc("/chat/group", handlers.HandleGroupMessages)

	port := ":" + os.Getenv("SERVER_PORT")

	log.Printf("Server is running on port %s", port[1:])
	log.Fatal(http.ListenAndServe(port, r))
}
