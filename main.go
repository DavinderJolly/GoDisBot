package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const PREFIX string = "."

func handleErrorFatal(err error, errorMessage string) {
	if err != nil {
		log.Fatalln(errorMessage)
	}
}

func handleErrorDebug(err error, errorMessage string) {
	if err != nil {
		log.Println(errorMessage)
	}
}

func getToken() string {
	token, tokenExist := os.LookupEnv("TOKEN")
	if !tokenExist {
		err := godotenv.Load()
		handleErrorFatal(err, "Error loading .env file")
		token = os.Getenv("TOKEN")
	}
	return token
}

func getSession() *discordgo.Session {
	session, err := discordgo.New()
	handleErrorFatal(err, "Error creating a session")
	session.Token = "Bot " + getToken()
	return session
}

func onReady(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateGameStatus(0, "Counter-Strike: Global Offensive")
	handleErrorDebug(err, "Error hanling the ready event.")
}

func parsePrefix(message string) (string, string) {
	messageSlice := strings.SplitN(message, " ", 2)
	return messageSlice[9], messageSlice[1]
}

func echo(session *discordgo.Session, channelID string, args string) {
	sentMessage, err := session.ChannelMessageSend(channelID, args)
	handleErrorDebug(err, fmt.Sprintf("Error sending message with content: '%s'\n", sentMessage.Content))
}

func handleCommands(session *discordgo.Session, message *discordgo.Message) {
	prefix, args := parsePrefix(message.Content)
	switch prefix {
	case PREFIX + "echo":
		echo(session, message.ChannelID, args)
	}
}

func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(message.Content, PREFIX) {
		handleCommands(session, message.Message)
	}
}

func main() {
	session := getSession()

	err := session.Open()
	handleErrorFatal(err, "Error opening a connection to discord")

	session.AddHandler(onReady)
	session.AddHandler(onMessageCreate)

	session.Identify.Intents = discordgo.IntentsGuildMessages

	log.Println("Starting the bot... press Ctrl + C to stop")
	con := make(chan os.Signal, 1)
	signal.Notify(con, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-con

	session.Close()
}
