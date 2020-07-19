package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg/botutil"
	dg "github.com/bwmarrin/discordgo"
	"strings"
)

func (bot *Bot) cmdAdd(msg *dg.Message, args []string) {
	var text, voice *dg.Channel

	// args should be [prefix, add, <text channel id>, <voice channel id>]
	if len(args) < 4 {
		util.Reply(bot.Client, msg,
			fmt.Sprintf("%s add <channel id> <channel id>", bot.Config.Prefix),
		)
		return
	}

	channelAID := args[2]
	if strings.HasPrefix(channelAID, "<") {
		channelAID = util.FilterTag(channelAID)
	}
	channelBID := args[3]
	if strings.HasPrefix(channelBID, "<") {
		channelBID = util.FilterTag(channelBID)
	}

	channelA, err := bot.Client.Channel(channelAID)

	if err != nil {
		util.Reply(bot.Client, msg, fmt.Sprintf(`"%s" isn't a channel ID`, args[2]))
		return
	}

	channelB, err := bot.Client.Channel(channelBID)

	if err != nil {
		util.Reply(bot.Client, msg, fmt.Sprintf(`"%s" isn't a channel ID`, args[3]))
		return
	}

	if channelA.Type == channelB.Type {
		util.Reply(bot.Client, msg, "Both channels must uniquely be a voice and text channel")
		return
	}

	if channelA.Type == dg.ChannelTypeGuildText {
		text = channelA
	} else if channelA.Type == dg.ChannelTypeGuildVoice {
		voice = channelA
	} else {
		util.Reply(bot.Client, msg,
			fmt.Sprintf(`"%s" isn't a voice or text channel.'`, channelAID),
		)
		return
	}

	if channelB.Type == dg.ChannelTypeGuildText {
		text = channelB
	} else if channelB.Type == dg.ChannelTypeGuildVoice {
		voice = channelB
	} else {
		util.Reply(bot.Client, msg,
			fmt.Sprintf(`"%s" isn't a voice or text channel.'`, channelBID),
		)
		return
	}

	if text != nil && voice != nil {
		bot.AddPair(text, voice)
		util.Reply(
			bot.Client,
			msg,
			fmt.Sprintf("Linked %s and %s", text.Name, voice.Name),
		)
	} else {
		util.Reply(bot.Client, msg, "An internal error occurred.")
	}
}

func (bot *Bot) cmdRemove(msg *dg.Message, args []string) {
	// args should be [prefix, remove, <channel id / #channel>]
	if len(args) < 3 {
		return
	}

	channel := util.FilterTag(args[2])
	removed := bot.RemPair(channel)

	if removed {
		util.Reply(bot.Client, msg, "Removed.")
	} else {
		util.Reply(bot.Client, msg, "Couldn't find "+channel)
	}
}

func (bot *Bot) cmdList(msg *dg.Message) {
	var text, voice *dg.Channel
	var err error
	var list = "Linked Channels\n"
	for voiceID, textID := range bot.Config.Channels {
		if voice, err = bot.Client.Channel(voiceID); err == nil {
			list += fmt.Sprintf(" - **%s** is linked with ", voice.Name)
		} else {
			list += fmt.Sprintf(" - Unknown (`%s`) is linked with ", voiceID)
		}
		if text, err = bot.Client.Channel(textID); err == nil {
			list += fmt.Sprintf("%s\n", text.Mention())
		} else {
			list += fmt.Sprintf("Unknown (`%s`)\n", textID)
		}
	}

	util.Reply(bot.Client, msg, list)
}
