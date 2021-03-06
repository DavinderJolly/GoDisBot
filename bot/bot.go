package bot

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// BotCommandPrefix is prefix for the bot's.
const BotCommandPrefix string = "."

func getToken() string {
	token, tokenExist := os.LookupEnv("TOKEN")
	if !tokenExist {
		err := godotenv.Load()
		HandleErrorFatal(err, "Error loading .env file")
		token = os.Getenv("TOKEN")
	}
	return token
}

func getSession() *discordgo.Session {
	session, err := discordgo.New()
	HandleErrorFatal(err, "Error creating a session")
	session.Token = "Bot " + getToken()
	return session
}

func onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(message.Content, BotCommandPrefix) {
		HandleCommands(session, message.Message)
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

// Run starts the bot
func Run() {
	session := getSession()

	err := session.Open()
	HandleErrorFatal(err, "Error opening a connection to discord")

	err = session.UpdateGameStatus(0, "Counter-Strike: Global Offensive")
	HandleErrorDebug(err, "Error setting bot status")

	setIntentsAndHandlers(session)

	log.Println("Starting the bot... press Ctrl + C to stop")

	// Keep running until interupted
	con := make(chan os.Signal, 1)
	signal.Notify(con, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-con
	log.Println("Closing...")
	session.Close()
}
