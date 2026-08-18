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

var users []net.Conn

func SendMessages(user net.Conn, input *bufio.Scanner) {
	for i := range users {
		if users[i] == user {
			continue
		}
		fmt.Fprintln(users[i], input.Text())
	}
}

var messageSize uint32 = 2048

func handleConn(conn net.Conn) {
	defer conn.Close()

	users = append(users, conn)

	var user user.User
	data := make([]byte, 200)

	n, err := conn.Read(data)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(data[:n], &user)
	if err != nil {
		log.Fatalln("possivel argumentos fora do range")
		log.Fatalln(err)
	}

	// json.Unmarshal([]byte(input.Text()), &user)
	fmt.Println("nova conexão no servidor:", user.Name)
	defer fmt.Printf("- %s foi embora\n", user.Name)

	input := bufio.NewScanner(conn)
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
