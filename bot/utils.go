package bot

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func ParsePrefix(message string) (string, string) {
	messageSlice := strings.SplitN(message, " ", 2)
	if len(messageSlice) < 2 {
		messageSlice = append(messageSlice, "")
	}
	return messageSlice[0], messageSlice[1]
}

func HasPermission(session *discordgo.Session, userID string, channelID string, perm int64) bool {
	userPerms, err := session.UserChannelPermissions(userID, channelID)
	HandleErrorDebug(err, "Error getting user channel permissions")
	return userPerms&perm == perm
}
