package NetCat

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	clients = make(map[net.Conn]string)
	mutex   = &sync.Mutex{}
)

// IsAllowedMessage checks if the message contains only allowed characters.
func sendPrompt(conn net.Conn, name string) {
	tm := time.Now().Format("2006-01-02 15:04:05")
	prompt := fmt.Sprintf("[%s] [%s]: ", tm, name)
	conn.Write([]byte(prompt))
}

// HandleConnection manages an individual client's connection to the chat server.
func HandleConnection(conn net.Conn) {
	defer conn.Close()
	welcomeLogo, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Println("error reading welcome file:", err)
		welcomeLogo = []byte("Welcome to the Chat!\n")
	}
	conn.Write(welcomeLogo)
	if err != nil {
		log.Println("error reading welcome file:", err)
		return
	}
	reader := bufio.NewReader(conn)
	var name string
	conn.Write([]byte("\n[ENTER YOUR NAME]:"))
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error connection Interruption", err)
			return
		}
		name = strings.TrimSpace(input)
		if name != "" && NameChecker(name) && LimitsChecker(conn) && LenChecker(name) && IsAllowedMessage(name) {
			break
		} else if !LimitsChecker(conn) {
			return
		} else if !LenChecker(name) {
			conn.Write([]byte("The Name u putted is bigger than 20 characters. Please choose another one.\n"))
		} else {
			conn.Write([]byte("This name is Invalid or already taken. Please choose another one.\n"))
		}
		conn.Write([]byte("[ENTER YOUR NAME]:"))
	}
	mutex.Lock()
	clients[conn] = name
	history, err := os.ReadFile("history.txt")
	if err != nil {
		log.Println("error reading history file:", err)
		history = []byte{}
	}
	if len(history) > 0 {
		conn.Write([]byte("\n--- Chat History ---\n"))
		conn.Write(history)
		conn.Write([]byte("--- End of History ---\n"))
	}
	mutex.Unlock()
	log.Printf("\n%s has joined our chat...", name)
	joinMsg := fmt.Sprintf("\n%s has joined our chat...\n", name)

	Broadcast(conn, joinMsg)
	for {
		sendPrompt(conn, name)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println(name, "left the channel. type of error : ", err)
			Broadcast(conn, fmt.Sprintf("\n%s has left the chat.\n", name))
			UserDelete(clients, conn)
			break
		}
		if !IsAllowedMessage(message) || strings.TrimSpace(message) == "" {
			continue
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		fnlmsg := fmt.Sprintf("[%s] [%s]: %s", tm, name, message)
		AppendToHistory(fnlmsg)
		fnlmsg = "\n" + fnlmsg
		Broadcast(conn, fnlmsg)
	}
}
