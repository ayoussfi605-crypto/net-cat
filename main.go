package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	port := "8989"
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(os.Args) == 2 {
		port = os.Args[1]
	}
	lisn, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Fatalf("error for listen: %v", err)
	}
	defer lisn.Close()
	log.Printf("listn in port : %s", port)

	for {
		conn, err := lisn.Accept()
		if err != nil {
			log.Printf("err for accept thisreader := bufio.NewReader connection: %v", err)
			continue
		}
		go welcommsg(conn)
	}
}

func welcommsg(conn net.Conn) {

	// Sdd l'connection dyal had l'client mli nsaliw mno
	defer conn.Close()
	// hlname(name)
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
	for{
		//read until find \n  
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Println("error connection Interruption", err)
		}
		// Nqy smiya mn spaces w l'retour à la ligne
		name = strings.TrimSpace(input)
		if name != ""{
			//name != "" left bocle 
			break
		}

		conn.Write([]byte("[ENTER YOUR NAME]:"))
		
	}
	fmt.Printf("%s has joined our chat...",name)
}

