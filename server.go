package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type geophone struct {
	id          string
	age         int
	currentRoom *room
	conn        net.Conn
}

type serverInfo struct {
	rooms map[string]*room
	lock  sync.Mutex
}

type room struct {
	name      string
	messages  []string
	geophones []*geophone
	lock      sync.Mutex
}

// postMessage takes a message from a client, and sends that message
// to all other clients in the same room
func (c *room) postMessage(message string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for _, receiver := range c.geophones {
		receiver.conn.Write([]byte(message))
	}
}

// addGeophone adds a geophone to a room so that other clients will be
// able to send them messages and other geophones can send them messages
func (c *room) addGeophone(u *geophone) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.geophones = append(c.geophones, u)
	u.currentRoom = c
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: server port\n")
		os.Exit(1)
	}
	// set up server
	service := ":" + os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	// initialize server data structures
	server := serverInfo{rooms: make(map[string]*room)}
	newRoom(&server, "all")

	// handle connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, &server)
	}
}

// promptUserStr writes a prompt to the connection and waits for a response,
// and returns that response as a string.
func promptUserStr(conn net.Conn, prompt string) (string, error) {
	conn.Write([]byte(prompt))
	response, err := bufio.NewReader(conn).ReadString('\n')
	response = strings.TrimSuffix(response, "\n")
	return response, err
}

// promptUserInt writes a prompt to the connection and waits for a response,
// and return that response as an int
func promptUserInt(conn net.Conn, prompt string) (int, error) {
	var response int

	conn.Write([]byte(prompt))
	_, err := fmt.Fscan(bufio.NewReader(conn), &response)
	return response, err
}

func handleClient(conn net.Conn, server *serverInfo) {
	defer conn.Close()
	daytime := time.Now().String()
	conn.Write(append([]byte(daytime), '\n'))

	// Get some information from the user
	//TODO: Get ID from the geophone
	id := "geophone"
	newGeophone := geophone{id: id, conn: conn}

	//TODO: Don't add all the geophone to the same room
	server.rooms["all"].addGeophone(&newGeophone)

	newGeophone.currentRoom.postMessage(id + " has entered the room.\n")
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == io.EOF {
			break
		}
		checkError(err)
		message = newGeophone.id + ": " + message
		newGeophone.currentRoom.postMessage(message)
	}
	newGeophone.currentRoom.postMessage(id + " has left the room.\n")
}

// newRoom creates a new room and adds it to the servers chat rooms
//  server: server struct
//  roomName: the name of the new chat room to be created
func newRoom(server *serverInfo, roomName string) {
	newRoom := room{name: roomName}
	server.rooms[newRoom.name] = &newRoom
}

func checkError(err error) {
	if err == io.EOF {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		os.Exit(1)
	}
}
