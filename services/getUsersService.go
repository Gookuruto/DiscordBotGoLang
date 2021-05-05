package services

import (
	"errors"
	"hourBot/database"
	"hourBot/database/models"
)

type UserService struct {
	Db *database.DatabaseConnection
}

func NewUserService(db *database.DatabaseConnection) *UserService {
	userService := new(UserService)

	userService.Db = db

	return userService
}

func (u *UserService) AddUser(userId, username string) error {
	user := models.User{UserId: userId, Username: username}
	existingUser := models.User{}
	u.Db.DB.Table("users").Find(&existingUser, "user_id = ?", userId)
	if existingUser.ID == 0 {
		u.Db.DB.Create(&user)
		return nil
	} else {
		return errors.New("User Already Exist")
	}
}

func (u *UserService) GetUser(userId string) *models.User {
	result := new(models.User)
	u.Db.DB.First(&result, "user_id = ?", userId)
	if result.ID == 0 {
		return nil
	}
	return result
}

func (u *UserService) GetAllUsers() *[]models.User {
	result := new([]models.User)
	u.Db.DB.Find(&result)
	if result == nil {
		return nil
	}
	return result
}
