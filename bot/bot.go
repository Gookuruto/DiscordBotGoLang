package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/diamondburned/arikawa/discord"
	"github.com/diamondburned/arikawa/gateway"
	"github.com/diamondburned/arikawa/session"
)

var bot *session.Session
var totalGuilds int
var guilds map[string]bool
var guildsBlackList []string
var bots map[string]bool
var hourBot *HourBot

func Run(dbConnection string) {
	bots = make(map[string]bool)
	guilds = make(map[string]bool)
	hourBot = NewHourBotInstance(dbConnection)
	//SETUP db connection

	//setup connection to discord
	session, err := session.New("Bot " + "TOKEN")

	session.Gateway.Identifier.Intents = 256 + 1 + 512 + 4096
	bot = session
	if err != nil {
		panic(err)
	}
	getUsers()
	err = session.Open()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.Gateway.UpdateStatus(gateway.UpdateStatusData{
		Game: &discord.Activity{
			Name: "@ to get stats",
		},
	})
	//Switched normal status and status displaying tracked server
	statusupdate := time.NewTicker(time.Second * 10)
	flip := false

	go func() {
		for {
			select {
			case <-statusupdate.C:
				var playingStr string
				if flip {
					playingStr = "Tracking stats for " + strconv.Itoa(totalGuilds) + " servers!"
					flip = false
				} else {
					playingStr = "@ to get stats"
					flip = true
				}
				session.Gateway.UpdateStatus(gateway.UpdateStatusData{
					Game: &discord.Activity{
						Name: playingStr,
					},
				})
			}
		}
	}()

	//adds handlers

	session.AddHandler(presenceUpdate)
	//session.AddHandler(guildAdded)
	session.AddHandler(newMessage)

	blacklists := os.Getenv("BLACKLIST")
	guildsBlackList = strings.Split(blacklists, ",")

	fmt.Println("Bot started ...")

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-exitChan

}

func getUsers() {
	users := hourBot.UserService.GetAllUsers()
	if users == nil {
		return
	}
	for _, v := range *users {
		DiscordUsers[v.UserId] = &discordUser{userID: v.UserId}
	}
}

func newMessage(msg *gateway.MessageCreateEvent) {
	if msg.Author.Bot {
		return
	}
	if msg.GuildID == 0 {
		fmt.Println("private message dont respond")
	} else {
		if msg.Content == "!track" {
			bot.SendMessage(msg.ChannelID, "Starting tracking your game time", nil)
			err := hourBot.UserService.AddUser(msg.Author.ID.String(), msg.Author.Username)
			if err == nil {
				DiscordUsers[msg.Author.ID.String()] = &discordUser{userID: msg.Author.ID.String()}
			}
		}
		if msg.Content == "!time" {
			bot.SendMessage(msg.ChannelID, "Creating your stats please wait...", nil)
			games := hourBot.TimeService.GetTopGames(msg.Author.ID.String(), 2)

			for _, v := range *games {
				hours := v.TimeMinutes / 60
				minutes := v.TimeMinutes % 60
				content := fmt.Sprintf("%s : %d hours %d minutes", v.Game.Name, hours, minutes)
				bot.SendMessage(msg.ChannelID, content, nil)
			}
		}
		//member, err := bot.Member(msg.GuildID, discord.UserID(mentionedUser))

	}
}

//Updating presence of user and start tracking games time there is a bug when we can disable and enable game this work with changing status to invisisble and back to online
func presenceUpdate(p *gateway.PresenceUpdateEvent) {
	if DiscordUsers[p.User.ID.String()] != nil && (p.Game == nil || p.Game.Type == discord.GameActivity) {
		if DiscordUsers[p.User.ID.String()].isPlaying && (p.Game == nil || DiscordUsers[p.User.ID.String()].currentGame != p.Presence.Game.Name) {
			hourBot.UpdateOrSave(p.User.ID.String(), DiscordUsers[p.User.ID.String()])
			DiscordUsers[p.User.ID.String()].reset()
			DiscordUsers[p.User.ID.String()].startTracking(p)
		} else {
			if p.Game == nil {
				hourBot.UpdateOrSave(p.User.ID.String(), DiscordUsers[p.User.ID.String()])
				DiscordUsers[p.User.ID.String()].reset()
			}
			if DiscordUsers[p.User.ID.String()].isPlaying == false && p.Game != nil {
				DiscordUsers[p.User.ID.String()].startTracking(p)
			}
		}
	}
}

// func handlePresenceUpdate(presence *discordgo.PresenceUpdate) {
// 	game := presence.Presence.
// 	user := discordUsers[presence.User.ID]
// 	user.mu.Lock()
// 	defer user.mu.Unlock()
// 	if game != nil { //Started Playing Game
// 		if game.Name != user.currentGame {
// 			fmt.Fprintln(out, "Started Playing Game "+game.Name)
// 			if user.isPlaying == true { //Switching from other game
// 				fmt.Fprintln(out, "Switching From Other Game "+user.currentGame)
// 				user.save()
// 				//saveGuild(user)
// 				user.reset()
// 				user.startTracking(presence)
// 			} else { //Not currently playing game
// 				fmt.Fprintln(out, "Not Playing Any Game")
// 				user.startTracking(presence)
// 			}
// 		}
// 	} else { //Stopped Playing Game
// 		if user.currentGame != "" {
// 			fmt.Fprintln(out, "Stopped Playing Game")
// 			user.save()
// 			//saveGuild(user)
// 			user.reset()
// 		}
// 	}
// }
