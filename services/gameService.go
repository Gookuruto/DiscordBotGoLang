package services

import (
	"hourBot/database"
	"hourBot/database/models"
)

type GameService struct {
	Db *database.DatabaseConnection
}

func NewGameService(db *database.DatabaseConnection) *GameService {
	gameService := new(GameService)

	gameService.Db = db

	return gameService
}

func (g *GameService) AddGame(gameName string) {
	game := models.Game{Name: gameName}
	existingGame := models.Game{}
	g.Db.DB.Table("games").Find(&existingGame, "name = ?", gameName)
	if existingGame.ID == 0 {
		g.Db.DB.Create(&game)
	}
}
