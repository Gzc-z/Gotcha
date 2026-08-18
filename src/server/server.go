// Package main
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"gochat/src/user"
)

var (
	users       []net.Conn
	nameSize    uint32 = 100
	messageSize uint32 = 2048
)

func SendMessages(user net.Conn, input *bufio.Scanner) {
	for i := range users {
		if users[i] == user {
			continue
		}
		fmt.Fprintln(users[i], input.Text())
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	users = append(users, conn)

	var user user.User

	input := bufio.NewScanner(conn)
	input.Buffer(
		make([]byte, nameSize),
		int(nameSize),
	)
	if !input.Scan() {
		if err := input.Err(); err != nil {
			log.Println("erro lendo conexão:", err)
		}
		return
	}
	// must accept valid structs like user
	// TODO: validate user
	err := json.Unmarshal([]byte(input.Text()), &user)
	if err != nil {
		fmt.Println(err)
	}

	// json.Unmarshal([]byte(input.Text()), &user)
	fmt.Println("nova conexão no servidor:", user.Name)
	defer fmt.Printf("- %s foi embora\n", user.Name)

	for input.Scan() {
		SendMessages(conn, input)
	}
}

// TODO do a whitelist, use dns and verify if a connection of ip contains in a whitelist
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
