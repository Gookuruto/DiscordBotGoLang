package bot

import (
	"hourBot/database"
	"hourBot/services"
	"time"
)

type HourBot struct {
	GameService *services.GameService
	TimeService *services.TimeService
	UserService *services.UserService
}

func (h *HourBot) UpdateOrSave(idToLookup string, user *discordUser) {
	currentUser := h.UserService.GetUser(idToLookup)
	if currentUser == nil {
		h.UserService.AddUser(user.userID, "")
	}
	h.GameService.AddGame(user.currentGame)
	timeMinutes := time.Now().Sub(user.startedPlaying).Minutes()
	h.TimeService.AddGameTime(user.currentGame, user.userID, int64(timeMinutes))

}

func NewHourBotInstance(connectionString string) *HourBot {
	db := database.NewDatabase(connectionString)
	hourBot := new(HourBot)
	hourBot.GameService = &services.GameService{Db: db}
	hourBot.TimeService = &services.TimeService{Db: db}
	hourBot.UserService = &services.UserService{Db: db}

	return hourBot
}
