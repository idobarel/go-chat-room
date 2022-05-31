package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	fmt.Println("new client has connected >> ", conn.RemoteAddr().String())
	c := &client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
	}
	c.readInput()
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room, 0),
		commands: make(chan command),
	}
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("Gotcha, %s!", c.nick))
}
func (s *server) join(c *client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c
	s.quitRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s Has Joined!", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s!", r.name))
}
func (s *server) listRooms(c *client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("Rooms >> %s", strings.Join(rooms, ", ")))
}
func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("You must join a room first!"))
		return
	}
	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
}

func (s *server) quit(c *client, args []string) {
	fmt.Printf("Client has disconnected >> '%s'.\n", c.conn.RemoteAddr().String())
	s.quitRoom(c)
	c.msg("Goodbye!")
	c.conn.Close()
}

func (s *server) quitRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s Has Left!", c.nick))
	}
}
