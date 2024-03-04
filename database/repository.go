package database

import (
	"github.com/AlexWilliam12/silent-signal/client"
	"github.com/AlexWilliam12/silent-signal/database/models"
)

// Create a user on database
func CreateUser(user client.Auth) (int64, error) {
	db := OpenConn()
	result := db.Select("username", "password").Create(&models.User{Username: user.Username, Password: user.Password})
	return result.RowsAffected, result.Error
}

// Find a single user on database
func FindUser(user client.Auth) (int64, error) {
	db := OpenConn()
	result := db.Where("users.username = ? AND users.password = ?", user.Username, user.Password).First(&models.User{Username: user.Username, Password: user.Password})
	return result.RowsAffected, result.Error
}

func FindContact(username string) (*models.User, error) {
	db := OpenConn()
	var user models.User
	result := db.Where("users.username = ?", username).Find(&user)
	return &user, result.Error
}

func DeleteUser(username string) (int64, error) {
	db := OpenConn()
	result := db.Unscoped().Where("users.username = ?", username).Delete(&models.User{})
	return result.RowsAffected, result.Error
}
