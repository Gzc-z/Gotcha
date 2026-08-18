// Package main
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"gochat/src/user"

	"github.com/google/uuid"
	"go.yaml.in/yaml/v3"
)

var userConfig = "src/userConf.yaml"

type Config struct {
	User *user.User
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Conn net.Conn
}

func getConfig() Config {
	data, err := os.ReadFile(userConfig)
	if err != nil {
		log.Fatalln(err)
	}
	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalln(err)
	}
	input := bufio.NewScanner(os.Stdin)

	// that will be executed once, then doesn't have a problem
	file, _ := os.Open("config.yaml")
	defer file.Close()

	config.User = &user.User{}
	if config.User.ID == uuid.Nil {
		config.User.ID = uuid.New()
		yaml.NewEncoder(file).Encode(config)
	}
	if config.User.Name == "" {
		fmt.Print("Seu nome: ")
		input.Scan()
		config.User.Name = input.Text()
		yaml.NewEncoder(file).Encode(config)
	}
	return config
}

func SendMessage(conf Config) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		msg := fmt.Sprintf("%s: %s", conf.User.Name, scanner.Text())
		fmt.Fprintln(conf.Conn, msg)
	}
	fmt.Fprintln(conf.Conn, conf.User.Name+" - disconected")
}

func main() {
	config := getConfig()
	conn, err := net.Dial("tcp4", config.Host+":"+config.Port)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	config.Conn = conn
	data, err := json.Marshal(config.User)
	if err != nil {
		log.Fatalln(err)
	}

	writer := bufio.NewWriter(conn)
	fmt.Fprintln(writer, string(data))
	if err := writer.Flush(); err != nil {
		log.Fatalln(err)
	}

	go SendMessage(config)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		fmt.Println(input.Text())
	}
}
