package internal

import (
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (bot *Bot) onMessage(_ *dg.Session, msg *dg.MessageCreate) {
	if msg.Author.Bot || !strings.HasPrefix(msg.Content, bot.Config.Prefix) {
		return
	}

	args := strings.Fields(msg.Content)

	// possible args
	// args = [prefix, add, <text channel id>, <voice channel id>]
	// args = [prefix, remove, <channel id>]
	if len(args) < 2 {
		return
	}

	isAdmin, _ := bot.Auth.Verify(msg.Author.ID)

	if isAdmin {
		switch args[1] {
		case "add":
			bot.cmdAdd(msg.Message, args)
			break
		case "remove":
			bot.cmdRemove(msg.Message, args)
			break
		case "list":
			bot.cmdList(msg.Message)
			break
		}
	} else {
		_, _ = util.Reply(
			bot.Client,
			msg.Message,
			"You need to be an admin to run these commands.",
		)
	}
}

func (bot *Bot) onReady(_ *dg.Session, ready *dg.Ready) {
	log.Printf(
		"Linked Channels - ready as %s#%s",
		ready.User.Username,
		ready.User.Discriminator,
	)
}

func (bot *Bot) onVoiceUpdate(_ *dg.Session, voice *dg.VoiceStateUpdate) {
	oldState, isOK := bot.OldStates[voice.UserID]
	bot.OldStates[voice.UserID] = voice.VoiceState

	if isOK {
		bot.handleOldState(oldState)
	}

	if textID, isOK := bot.Config.Channels[voice.ChannelID]; isOK {
		text, err := bot.Client.Channel(textID)
		voice, _ := bot.Client.Channel(voice.ChannelID)
		if err == nil {
			bot.sync(text, voice)
		} else {
			log.Printf("Couldn't find %s linked with %s\n", textID, voice.Name)
		}
	}
}

func (bot *Bot) handleOldState(voice *dg.VoiceState) {
	if textID, isOK := bot.Config.Channels[voice.ChannelID]; isOK {
		text, err := bot.Client.Channel(textID)
		voice, _ := bot.Client.Channel(voice.ChannelID)
		if err == nil {
			bot.sync(text, voice)
		} else {
			log.Printf("Couldn't find %s linked with %s\n", textID, voice.Name)
		}
	}
}
