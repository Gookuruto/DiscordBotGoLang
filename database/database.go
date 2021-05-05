package database

import (
	"database/sql"
	"hourBot/database/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DbConnection interface {
	ConnectDatabase(connectionString string)
}
type DatabaseConnection struct {
	DB *gorm.DB
}

func NewDatabase(connectionString string) *DatabaseConnection {
	db := new(DatabaseConnection)
	db.DB = new(gorm.DB)
	db.ConectDatabase(connectionString)

	return db
}

func (db *DatabaseConnection) ConectDatabase(connectionString string) {
	database, err := gorm.Open("postgres", connectionString)
	if err != nil {
		createPgDb()
		database, err = gorm.Open("postgres", connectionString)
		if err != nil {
			panic(err)
		}
	}
	database.AutoMigrate(&models.Game{})
	database.AutoMigrate(&models.GameTime{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT").AddForeignKey("game_id", "games(id)", "RESTRICT", "RESTRICT")
	database.AutoMigrate(&models.User{})

	db.DB = database
}

func createPgDb() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=pass sslmode=disable")
	if err != nil {
		panic(err)
	}
	dbName := "gorm"
	_, err = db.Exec("create database " + dbName)
	if err != nil {
		//handle the error
		log.Fatal(err)
	}

}
