package models

type GameTime struct {
	ID          int64
	GameID      int64
	Game        Game
	TimeMinutes int64
	UserID      int64
	User        User
}
