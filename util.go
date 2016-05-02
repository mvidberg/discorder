package main

import (
	"github.com/jonas747/discordgo"
	"log"
)

// For logs, should probably move this somewhere else though
func (app *App) Write(p []byte) (n int, err error) {
	cop := string(p)
	app.HandleLogMessage(cop)
	return len(p), nil
}

func (app *App) GetNotificationSettingsForChannel(channelId string) *ChannelNotificationSettings {
	channel, err := app.session.Channel(channelId)
	if err != nil {
		log.Println("Error getting channel from state", err)
		return nil
	}

	if channel.IsPrivate {
		return &ChannelNotificationSettings{Notifications: ChannelNotificationsAll}
	}

	for _, gs := range app.guildSettings {
		if gs.GuildID == channel.GuildID {

			cn := &ChannelNotificationSettings{
				Notifications:    gs.MessageNotifications,
				Muted:            gs.Muted,
				SurpressEveryone: gs.SupressEveryone,
			}
			if gs.Muted {
				return cn
			}
			for _, override := range gs.ChannelOverrides {
				if override.ChannelID == channel.ID {
					cn.Notifications = override.MessageNotifications
					cn.Muted = override.Muted
					break
				}
			}
			return cn
		}
	}

	// Use default guild settings
	guild, err := app.session.Guild(channel.GuildID)
	if err != nil {
		log.Println("Error getting guild from state", err)
		return nil
	}
	return &ChannelNotificationSettings{
		Notifications: guild.DefaultMessageNotifications,
	}
}

// Compare readstate's last_message to channel's last_message and if theres new show so
// Also number of mentions
// Take notifications settings into mind also
func (app *App) GetStartNotifications() {
	// readStates := app.session.State.ReadState
	// for _, state := range readStates {
	// }
}

const (
	ChannelNotificationsAll      = 0
	ChannelNotificationsMentions = 1
	ChannelNotificationsNothing  = 2
)

type ChannelNotificationSettings struct {
	Notifications    int // 0 all, 1 mentions, 2 nothing
	Muted            bool
	SurpressEveryone bool
}

func GetChannelNameOrRecipient(channel *discordgo.Channel) string {
	if channel.IsPrivate {
		if channel.Recipient != nil {
			return channel.Recipient.Username
		} else {
			return "Recipient is nil!?"
		}
	}
	return channel.Name
}

func GetMessageAuthor(msg *discordgo.Message) string {
	if msg.Author != nil {
		return msg.Author.Username
	}
	return "Unknwon?"
}

func (app *App) IsFirstChannelMessage(channelId, msgId string) bool {
	first, ok := app.firstMessages[channelId]
	if !ok {
		return false
	}

	if first == msgId {
		return true
	}
	return false
}
