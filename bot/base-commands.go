package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func HandleCommands(session *discordgo.Session, message *discordgo.Message) {
	messagePrefix, args := ParsePrefix(message.Content)
	switch strings.ToLower(messagePrefix) {
	case BOT_COMMAND_PREFIX + "echo":
		Echo(session, message.ChannelID, args)
	case BOT_COMMAND_PREFIX + "ping":
		Ping(session, message.ChannelID, message.Timestamp)
	case BOT_COMMAND_PREFIX + "kick":
		Kick(session, message)
	case BOT_COMMAND_PREFIX + "tempban":
		TempBan(session, message)
	case BOT_COMMAND_PREFIX + "unban":
		Unban(session, message)
	case BOT_COMMAND_PREFIX + "purge":
		Purge(session, message)
	default:
		log.Println("Command not found")
	}
}

func Echo(session *discordgo.Session, channelID string, args string) {
	sentMessage, err := session.ChannelMessageSend(channelID, args)
	HandleErrorDebug(err, "Error sending message with content: '"+sentMessage.Content+"'")
}

func Ping(session *discordgo.Session, channelID string, timestamp discordgo.Timestamp) {
	pingTime, err := timestamp.Parse()
	HandleErrorDebug(err, "Could not parse timestamp")
	pingStr := strconv.Itoa(int(time.Now().UnixMilli()-pingTime.UnixMilli())) + "ms"
	sentMessage, err := session.ChannelMessageSend(channelID, pingStr)
	HandleErrorDebug(err, "Error sending message with content: '"+sentMessage.Content+"'")
}
