package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func getToken() string {
	token, tokenExist := os.LookupEnv("TOKEN")
	if !tokenExist {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln("Error loading .env file")
		}
	} else {
		token = os.Getenv("TOKEN")
	}
	return token
}

func main() {
	session := getSession()

	err := session.Open()
	if err != nil {
		log.Fatalln("Error opening a connection to discord")
	}

	log.Println("Starting the bot... press Ctrl + C to stop")
	con := make(chan os.Signal, 1)
	signal.Notify(con, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-con

	session.Close()
}

func getSession() *discordgo.Session {
	session, err := discordgo.New()
	if err != nil {
		log.Fatalln("Error creating a session")
	}
	session.Token = getToken()
	return session
}
