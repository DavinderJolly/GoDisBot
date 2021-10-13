package bot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// ParsePrefix returns the prefix + command and rest of the message split up
func ParsePrefix(message string) (string, string) {
	messageSlice := strings.SplitN(message, " ", 2)
	if len(messageSlice) < 2 {
		messageSlice = append(messageSlice, "")
	}
	return messageSlice[0], messageSlice[1]
}

// HasPermission checks if the user given has the given permission
func HasPermission(session *discordgo.Session, userID string, channelID string, perm int64) bool {
	userPerms, err := session.UserChannelPermissions(userID, channelID)
	HandleErrorDebug(err, "Error getting user channel permissions")
	return userPerms&perm == perm
}

// HandleErrorFatal checks if there is an error and stops the program with an error
func HandleErrorFatal(err error, errorMessage string) {
	if err != nil {
		log.Fatalln(errorMessage)
	}
}

// HandleErrorDebug checks if there is an error and logs it
func HandleErrorDebug(err error, errorMessage string) {
	if err != nil {
		log.Println(errorMessage)
	}
}
