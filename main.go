package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const botCommandPrefix string = "."

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
	log.Println("Bot is Ready")
	err := session.UpdateGameStatus(0, "Counter-Strike: Global Offensive")
	handleErrorDebug(err, "Error hanling the ready event.")
}

func parsePrefix(message string) (string, string) {
	messageSlice := strings.SplitN(message, " ", 2)
	if len(messageSlice) < 2 {
		messageSlice = append(messageSlice, "")
	}
	return messageSlice[0], messageSlice[1]
}

func echo(session *discordgo.Session, channelID string, args string) {
	sentMessage, err := session.ChannelMessageSend(channelID, args)
	handleErrorDebug(err, "Error sending message with content: '"+sentMessage.Content+"'")
}

func ping(session *discordgo.Session, channelID string) {
	sentMessage, err := session.ChannelMessageSend(channelID, session.HeartbeatLatency().String())
	handleErrorDebug(err, "Error sending message with content: '"+sentMessage.Content+"'")
}

func handleCommands(session *discordgo.Session, message *discordgo.Message) {
	messagePrefix, args := parsePrefix(message.Content)
	switch strings.ToLower(messagePrefix) {
	case botCommandPrefix + "echo":
		echo(session, message.ChannelID, args)
	case botCommandPrefix + "ping":
		ping(session, message.ChannelID)
	default:
		log.Println("Command not found")
	}
}

func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(message.Content, botCommandPrefix) {
		handleCommands(session, message.Message)
	}
}

func setIntentsAndHandlers(session *discordgo.Session) {
	session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageTyping |
		discordgo.IntentsDirectMessages |
		discordgo.IntentsDirectMessageTyping |
		discordgo.IntentsGuildPresences

	session.AddHandler(onMessageCreate)
}

func main() {
	session := getSession()

	err := session.Open()
	handleErrorFatal(err, "Error opening a connection to discord")

	err = session.UpdateGameStatus(0, "Counter-Strike: Global Offensive")
	handleErrorDebug(err, "Error hanling the ready event.")

	setIntentsAndHandlers(session)

	log.Println("Starting the bot... press Ctrl + C to stop")

	// Keep running until interupted
	con := make(chan os.Signal, 1)
	signal.Notify(con, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-con
	log.Println("Closing...")
	session.Close()
}
