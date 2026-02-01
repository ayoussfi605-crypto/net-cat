package main

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

func HandleConnection(conn net.Conn) {
	// Sdd l'connection dyal had l'client mli nsaliw mno
	// defer conn.Close()
	welcomeLogo, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Println("error reading welcome file:", err)
		return
	}
	conn.Write(welcomeLogo)
	// wahd reader smaile kaireadi smiya dyal l'client
	reader := bufio.NewReader(conn)
	var name string
	conn.Write([]byte("[ENTER YOUR NAME]:"))
	for {
		// read until find \n
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error connection Interruption", err)
			return
		}
		// Nqy smiya mn spaces w l'retour à la ligne
		name = strings.TrimSpace(input)
		if name != "" {
			// name != "" left bocle
			break
		}

		conn.Write([]byte("[ENTER YOUR NAME]:"))
	}
	mutex.Lock()
	clients[conn] = name
	mutex.Unlock()
	log.Printf("%s has joined our chat... \n", name)
	joinMsg := fmt.Sprintf("%s has joined our chat...\n", name)

	Broadcast(conn, joinMsg)
// ----------------------------------
	tm := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] [%s]",tm, name)
    conn.Write([]byte(msg))
	for{
		message, err := reader.ReadString('\n')
		if err != nil{
		log.Println("error connection Interruption", err)
			break
		}
		fnlmsg := fmt.Sprintf("%s %s %s", msg, name, message)
		Broadcast(conn, fnlmsg)
	}

}