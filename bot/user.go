package bot

import (
	"sync"
	"time"

	"github.com/diamondburned/arikawa/gateway"
)

var DiscordUsers = make(map[string]*discordUser)

type discordUser struct {
	userID         string
	mainGuild      string
	otherGuilds    map[string]string
	currentGame    string
	isPlaying      bool
	startedPlaying time.Time
	mu             sync.Mutex
}

// func (user *discordUser) save() {
// 	updateOrSave(user.userID, user)
// }

// func saveGuild(user *discordUser) {
// 	updateOrSave(user.mainGuild, user)
// 	for _, item := range user.otherGuilds {
// 		updateOrSave(item, user)
// 	}
// }

func (user *discordUser) startTracking(presence *gateway.PresenceUpdateEvent) {
	if presence.Game == nil {
		user.currentGame = ""
		user.isPlaying = false
		user.startedPlaying = time.Time{}
	} else {
		user.currentGame = presence.Game.Name
		user.isPlaying = true
		if presence.Game.Timestamps.Start != 0 {
			//TODO Make it work that disabling and enabling vs code discord will not double this time
			user.startedPlaying = time.Unix(int64(presence.Game.Timestamps.Start)/1000, 0)
		} else {
			user.startedPlaying = time.Now()
		}
	}
}

func (user *discordUser) reset() {
	user.isPlaying = false
	user.startedPlaying = time.Time{}
	user.currentGame = ""
}
