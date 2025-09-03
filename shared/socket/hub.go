package socket

import (
	"fmt"
	"log"
	"strings"
	"sync"

	socketio "github.com/googollee/go-socket.io"
)

var (
	userSockets = make(map[string]socketio.Conn)
	roomUsers   = make(map[string]map[string]bool) // room -> userId set

	serverInstance *socketio.Server
	once           sync.Once
)

// GetSocketIOServer returns the singleton Socket.IO server
func GetSocketIOServer() *socketio.Server {
	once.Do(func() {
		serverInstance = socketio.NewServer(nil)
		initSocketEvents(serverInstance)
	})
	return serverInstance
}
func CloseSocketIOServer() {
	if serverInstance != nil {
		log.Println("Shutting down Socket.IO server...")

		// Disconnect all connected users
		for userId, conn := range userSockets {
			log.Printf("Disconnecting user %s", userId)
			conn.Close()
		}

		// Close the server
		serverInstance.Close()
	}
}

func initSocketEvents(server *socketio.Server) {
	server.OnConnect("/", func(s socketio.Conn) error {
		u := s.URL()
		userId := strings.TrimSuffix(u.Query().Get("userId"), "/")
		fmt.Println("userId", userId)

		if userId == "" {
			log.Printf("connected: %s (no userId)", s.ID())
		} else {
			userSockets[userId] = s
			log.Printf("connected user %s with socket %s\n", userId, s.ID())
		}
		return nil 
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
		for uid, conn := range userSockets {
			if conn.ID() == s.ID() {
				delete(userSockets, uid)
				log.Printf("user %s disconnected\n", uid)

				// Remove user from all rooms
				for room, users := range roomUsers {
					if users[uid] {
						delete(users, uid)
						onlineList := make([]string, 0, len(users))
						for u := range users {
							onlineList = append(onlineList, u)
						}
						server.BroadcastToRoom("/", room, "onlineUsers", onlineList)
					}
				}
				break
			}
		}
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, data map[string]string) {
		room := data["room"]
		userId := data["userId"]
		if room == "" || userId == "" {
			log.Println("empty room or userId, ignoring join")
			return
		}

		s.Join(room)
		log.Printf("%s (user %s) joined room %s\n", s.ID(), userId, room)

		if _, ok := roomUsers[room]; !ok {
			roomUsers[room] = make(map[string]bool)
		}
		roomUsers[room][userId] = true

		onlineList := make([]string, 0, len(roomUsers[room]))
		for uid := range roomUsers[room] {
			onlineList = append(onlineList, uid)
		}
		server.BroadcastToRoom("/", room, "onlineUsers", onlineList)
	})

	server.OnEvent("/", "leaveRoom", func(s socketio.Conn, data map[string]string) {
		room := data["room"]
		userId := data["userId"]
		if room == "" || userId == "" {
			return
		}

		s.Leave(room)
		log.Printf("%s (user %s) left room %s\n", s.ID(), userId, room)

		if users, ok := roomUsers[room]; ok {
			delete(users, userId)
			if len(users) == 0 {
				delete(roomUsers, room)
			} else {
				onlineList := make([]string, 0, len(users))
				for uid := range users {
					onlineList = append(onlineList, uid)
				}
				server.BroadcastToRoom("/", room, "onlineUsers", onlineList)
			}
		}
	})

	server.OnEvent("/", "sendMessage", func(s socketio.Conn, msg map[string]string) {
		room := msg["room"]
		message := msg["message"]
		fmt.Println(message, "message")
		log.Printf("[%s] %s: %s\n", room, s.ID(), message)

		server.BroadcastToRoom("/", room, "newMessage", map[string]string{
			"from":    s.ID(),
			"message": message,
		})
	})
}
