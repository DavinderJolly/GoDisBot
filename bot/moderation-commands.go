package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Kick command kicks a user with a reason
func Kick(session *discordgo.Session, message *discordgo.Message) {
	if HasPermission(session, message.Author.ID, message.ChannelID, discordgo.PermissionKickMembers) {
		_, reason := ParsePrefix(message.ContentWithMentionsReplaced())
		for _, mention := range message.Mentions {
			err := session.GuildMemberDeleteWithReason(message.GuildID, mention.ID, reason)
			HandleErrorDebug(err, "Error while kicking the user")
			_, err = session.ChannelMessageSend(message.ChannelID, mention.Username+" kicked for reason: "+reason)
			HandleErrorDebug(err, "Error sending kick message for user "+mention.Username)
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "you don't have permission to use this command")
		HandleErrorDebug(err, "error sending not enough permission message")
	}
}

// TempBan temporarily bans a user with a reason
func TempBan(session *discordgo.Session, message *discordgo.Message) {
	if HasPermission(session, message.Author.ID, message.ChannelID, discordgo.PermissionBanMembers) {
		_, args := ParsePrefix(message.ContentWithMentionsReplaced())
		daysAndReason := strings.SplitN(args, " ", 2)
		days, err := strconv.Atoi(daysAndReason[0])
		HandleErrorDebug(err, "Error while splitting days and reason in tempBan command")
		reason := daysAndReason[1]
		for _, mention := range message.Mentions {
			session.GuildBanCreateWithReason(message.GuildID, mention.ID, reason, days)
			session.GuildMemberDeleteWithReason(message.GuildID, mention.ID, reason)
			_, err := session.ChannelMessageSend(message.ChannelID, mention.Username+" banned for reason: "+reason)
			HandleErrorDebug(err, "Error while kicking the user")
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "you don't have permission to use this command")
		HandleErrorDebug(err, "error sending not enough permission message")
	}
}

// Unban removes the ban on a user
func Unban(session *discordgo.Session, message *discordgo.Message) {
	if HasPermission(session, message.Author.ID, message.ChannelID, discordgo.PermissionBanMembers) {
		name := strings.ToLower(strings.SplitN(message.Content, " ", 3)[1])
		bans, err := session.GuildBans(message.GuildID)
		HandleErrorDebug(err, "Error getting bans from guild")
		for _, ban := range bans {
			if strings.ToLower(ban.User.Username) == name {
				err := session.GuildBanDelete(message.GuildID, ban.User.ID)
				HandleErrorDebug(err, "Error removing ban on "+ban.User.Username)
				_, err = session.ChannelMessageSend(message.ChannelID, ban.User.Username+" unbanned")
				HandleErrorDebug(err, "Error sending message after unban")
			}
		}
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "you don't have permission to use this command")
		HandleErrorDebug(err, "error sending not enough permission message")
	}
}

// Purge command removes said number of messages from the channel its used in
func Purge(session *discordgo.Session, message *discordgo.Message) {
	if HasPermission(session, message.Author.ID, message.ChannelID, discordgo.PermissionManageMessages) {
		splitMessage := strings.SplitN(message.ContentWithMentionsReplaced(), " ", 3)
		num, err := strconv.Atoi(splitMessage[1])
		HandleErrorDebug(err, "Error on purge command given number not integer")

		channelMessages, err := session.ChannelMessages(message.ChannelID, num, message.ID, "", "")
		HandleErrorDebug(err, "Error getting channel messages in purge")
		for _, channelMessage := range channelMessages {
			err = session.ChannelMessageDelete(message.ChannelID, channelMessage.ID)
			log.Println(err)
			HandleErrorDebug(err, "Error deleting messages with purge command on channel with id: "+message.ChannelID)
		}
		err = session.ChannelMessageDelete(message.ChannelID, message.ID)
		HandleErrorDebug(err, "Error deleting author's message in purge command")
	} else {
		_, err := session.ChannelMessageSend(message.ChannelID, "you don't have permission to use this command")
		HandleErrorDebug(err, "error sending not enough permission message")
	}
}
