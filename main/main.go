package main

import (
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

func main() {
	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		//s.Join("common")
		log.Println("connected:", s.ID())

		return nil
	})

	server.OnConnect("/draw", func(s socketio.Conn) error {
		s.SetContext("")
		s.Join("common")
		log.Println("connected /draw:", s.ID())

		return nil
	})

	server.OnEvent("/draw", "msg", func(s socketio.Conn, msg string) string {
		//s.SetContext(msg)

		fmt.Println("browacasting")
		server.ForEach("/draw", "common", func(conn socketio.Conn) {
			if conn.ID() == s.ID() {
				fmt.Println("same connection")
				return
			}
			conn.Emit("stroke", msg)
		})
		//server.BroadcastToRoom("/draw", "common", "stroke", msg)
		//fmt.Println("browadcast done", done)
		return "recv "
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("../client")))

	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
