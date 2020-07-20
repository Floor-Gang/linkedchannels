package internal

import (
	auth "github.com/Floor-Gang/authclient"
	dg "github.com/bwmarrin/discordgo"
	"log"
)

type Bot struct {
	Auth          *auth.AuthClient
	Client        *dg.Session
	Config        *Config
	OldStates     map[string]*dg.VoiceState
	Serving       string
	PermReference int
}

func Start() {
	config := GetConfig()
	// Start auth server
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to auth server", err)
	}

	register, err := authClient.Register(
		auth.Feature{
			Name:        "Linked Channels",
			Description: "Hide text channels until users enter the corresponding voice channel.",
			Commands: []auth.SubCommand{
				{
					Name:        "add",
					Description: "Add a new channel pair",
					Example:     []string{"add", "text channel ID", "voice channel ID"},
				},
				{
					Name:        "remove",
					Description: "Remove a channel pair",
					Example:     []string{"remove", "channel ID"},
				},
				{
					Name:        "list",
					Description: "List all the linked channels",
					Example:     []string{"list"},
				},
			},
			CommandPrefix: config.Prefix,
		})

	if err != nil {
		log.Fatalln("Failed to register with authentication server", err)
	}

	// setup discord bot
	client, _ := dg.New(register.Token)

	client.Identify.Intents = dg.MakeIntent(
		dg.IntentsGuildVoiceStates + dg.IntentsGuildMessages,
	)

	// setup bot struct
	bot := Bot{
		Auth:          &authClient,
		Client:        client,
		Config:        &config,
		Serving:       register.Serving,
		OldStates:     make(map[string]*dg.VoiceState),
		PermReference: dg.PermissionViewChannel,
	}

	client.AddHandler(bot.onMessage)
	client.AddHandler(bot.onVoiceUpdate)
	client.AddHandlerOnce(bot.onReady)

	if err = client.Open(); err != nil {
		log.Fatalln("Failed to connect to Discord", err)
	}
}

func (bot *Bot) AddPair(text *dg.Channel, voice *dg.Channel) {
	if voice.Type != dg.ChannelTypeGuildVoice {
		log.Fatalln("The channel provided isn't a voice channel.")
	} else if text.Type != dg.ChannelTypeGuildText {
		log.Fatalln("The channel provided isn't a text channel.")
	}

	bot.Config.Channels[voice.ID] = text.ID
	bot.Config.Save()
}

func (bot *Bot) RemPair(ID string) bool {
	for textID, voiceID := range bot.Config.Channels {
		if textID == ID || voiceID == ID {
			delete(bot.Config.Channels, textID)
			bot.Config.Save()
			return true
		}
	}
	return false
}

func (bot *Bot) getMembersOfVC(voiceID string) []string {
	var members []string
	guild, err := bot.Client.State.Guild(bot.Serving)

	if err != nil {
		log.Println("Failed to get guild " + bot.Serving)
		return members
	}

	for _, voiceState := range guild.VoiceStates {
		if voiceState.ChannelID == voiceID {
			members = append(members, voiceState.UserID)
		}
	}

	return members
}

func (bot *Bot) sync(text *dg.Channel, voice *dg.Channel) {
	var err error
	var inVC bool
	members := bot.getMembersOfVC(voice.ID)
	toAdd := members

	log.Printf("Syncing %s and #%s\n", voice.Name, text.Name)

	for _, perm := range text.PermissionOverwrites {
		if perm.Type == "member" {
			inVC = isInVC(perm.ID, members)
			if !inVC {
				log.Println("Removed " + perm.ID)
				err = bot.Client.ChannelPermissionDelete(text.ID, perm.ID)
				if err != nil {
					log.Printf("Failed to remove %s because \n"+err.Error(), perm.ID)
				}
			} else {
				log.Printf("%s is already in the voice channel\n", perm.ID)
				for i, memberID := range members {
					if memberID == perm.ID {
						members[i] = members[len(members)-1]
						members = members[:len(members)-1]
					}
				}
			}
		}

	}
	for _, memberID := range toAdd {
		log.Println("Added " + memberID)

		if err = bot.Client.ChannelPermissionSet(
			text.ID,
			memberID,
			"member",
			bot.PermReference,
			0,
		); err != nil {
			log.Printf("Failed to add %s because\n"+err.Error(), memberID)
		}
	}
}

func isInVC(memberID string, members []string) bool {
	for _, member := range members {
		if memberID == member {
			return true
		}
	}
	return false
}
