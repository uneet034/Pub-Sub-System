package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"pubsub/pubsub"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func autoId() (string) {

	return uuid.Must(uuid.NewUUID()).String()
}

var ps = &pubsub.PubSub{}

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {

		return true

	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}



	user := pubsub.User{
		Id:   autoId(),
		Connection: conn,
	}
	// add this user into the list

	ps.AddUser(user)

	fmt.Println("New User is connected, total:", len(ps.Users))


	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Something went wrong", err)

			// this user is disconnect or error connection we do need to remove subscriptions and remove user from pubsub

			ps.RemoveUser(user)
			log.Println("total users and subscriptions ", len(ps.Users), len(ps.Subscriptions))

			return
		}


		ps.Ack(user, messageType, p)

	}

}


func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "html")

	})

	http.HandleFunc("/ws", websocketHandler)


	http.ListenAndServe(":1000", nil)

	fmt.Println("Server is running: http://localhost:1000")

}