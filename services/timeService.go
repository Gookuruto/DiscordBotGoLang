package services

import (
	"hourBot/database"
	"hourBot/database/models"
)

type TimeService struct {
	Db *database.DatabaseConnection
}

func NewTimeService(db *database.DatabaseConnection) *TimeService {
	timeService := new(TimeService)

	timeService.Db = db

	return timeService
}

func (g *TimeService) AddGameTime(gameName, userId string, minutes int64) {
	user := new(models.User)
	g.Db.DB.First(user, "user_id = ?", userId)
	game := new(models.Game)
	g.Db.DB.First(game, "name = ?", gameName)
	gameTime := models.GameTime{UserID: user.ID, GameID: game.ID, TimeMinutes: minutes}
	existingGameTime := models.GameTime{}
	g.Db.DB.Find(&existingGameTime, "user_id = ? AND game_id = ?", user.ID, game.ID)
	if existingGameTime.ID == 0 {
		g.Db.DB.Create(&gameTime)
	} else {
		g.Db.DB.Model(&existingGameTime).Update("time_minutes", existingGameTime.TimeMinutes+minutes)
	}
}

func (g *TimeService) GetTimesForUser(userId string) *[]models.GameTime {
	gameTimes := new([]models.GameTime)

	g.Db.DB.Preload("User", "user_id = ?", userId).Preload("Game").Find(&gameTimes)

	return gameTimes
}

func (g *TimeService) GetTopGames(userId string, numberOfGames int) *[]models.GameTime {
	gameTimes := new([]models.GameTime)

	g.Db.DB.Preload("User", "user_id = ?", userId).Preload("Game").Order("time_minutes desc").Limit(numberOfGames).Find(&gameTimes)
	return gameTimes
}
