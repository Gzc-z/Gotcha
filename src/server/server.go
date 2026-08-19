// Package main
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"

	"gochat/src/user"

	"github.com/google/uuid"
)

var (
	users       []net.Conn
	nameSize    uint32 = 100
	messageSize uint32 = 2048
)

func sendMessages(user net.Conn, input *bufio.Scanner) {
	for i := range users {
		if users[i] == user {
			continue
		}
		fmt.Fprintln(users[i], input.Text())
	}
}

func httpHandler(conn net.Conn) {
	// oh gosh
	msg := "não podi se conectar via navegador"
	response := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
		"Content-Type: text/plain; charset=utf-8\r\n"+
		"Content-Length: %d\r\n"+
		"\r\n%s\n",
		len(msg),
		msg,
	)
	conn.Write([]byte(response))
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	users = append(users, conn)

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(
		make([]byte, nameSize),
		int(nameSize),
	)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			log.Println("erro lendo scanners:", err)
		}
		return
	}
	userMetadata := scanner.Text() // i don't think i should use this

	if strings.Contains(scanner.Text(), "GET /") {
		httpHandler(conn)
		return
	}

	var user user.User
	err := json.Unmarshal([]byte(userMetadata), &user)
	if err != nil {
		user.Name = "default-user-" + uuid.New().String()[0:8]
		user.ID = uuid.New()
		msg := "WARN: don't find user metadata, using default user"
		fmt.Printf("\033[1;33m%s\033[0m", msg)
	}
	defer fmt.Printf("- %s foi embora\n", user.Name)

	fmt.Println("nova conexão no servidor:", user.Name)

	for scanner.Scan() {
		sendMessages(conn, scanner)
	}
}

func main() {
	PORT := 3000

	ln, err := net.Listen("tcp4", fmt.Sprint(":", PORT))
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Println("listening")

	for {
		// TODO security and errors
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConn(conn)
	}
}
