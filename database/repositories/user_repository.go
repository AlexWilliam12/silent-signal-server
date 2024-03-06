package repositories

import (
	"github.com/AlexWilliam12/silent-signal/client"
	"github.com/AlexWilliam12/silent-signal/database"
	"github.com/AlexWilliam12/silent-signal/database/models"
)

// Create a user on database
func CreateUser(user client.UserRequest) (int64, error) {
	db := database.OpenConn()
	result := db.Select("username", "password").Create(&models.User{Username: user.Username, Password: user.Password})
	return result.RowsAffected, result.Error
}

// Find a single user on database querying by credentials
func FindUserByCredentials(user client.UserRequest) (*models.User, error) {
	db := database.OpenConn()
	var userQueried models.User
	result := db.Where("users.username = ? AND users.password = ?", user.Username, user.Password).First(&userQueried)
	return &userQueried, result.Error
}

// Find a single user on database querying by username
func FindUserByName(username string) (*models.User, error) {
	db := database.OpenConn()
	var userQueried models.User
	result := db.Where("users.username = ?", username).First(&userQueried)
	return &userQueried, result.Error
}

// Update user credentials
func UpdateUser(user *models.User) (int64, error) {
	db := database.OpenConn()
	result := db.Save(user)
	return result.RowsAffected, result.Error
}

// Delete user on database where name matches
func DeleteUserByName(username string) (int64, error) {
	db := database.OpenConn()
	result := db.Unscoped().Where("users.username = ?", username).Delete(&models.User{})
	return result.RowsAffected, result.Error
}
