package repositories

import (
	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
)

func CreateUser(user client.UserRequest) (int64, error) {
	db := database.OpenConn()
	result := db.Select("email", "username", "password").Create(&models.User{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	})
	return result.RowsAffected, result.Error
}

func FetchAllByUserames(username []string) ([]*models.User, error) {
	db := database.OpenConn()
	var users []*models.User
	result := db.Where("users.username IN ?", username).Find(users)
	return users, result.Error
}

func FindUserByCredentials(user client.UserRequest) (*models.User, error) {
	db := database.OpenConn()
	var persistedUser models.User
	result := db.Where("users.username = ? AND users.password = ?", user.Username, user.Password).First(&persistedUser)
	return &persistedUser, result.Error
}

func FindUserByName(username string) (*models.User, error) {
	db := database.OpenConn()
	var userQueried models.User
	result := db.Where("users.username = ?", username).First(&userQueried)
	return &userQueried, result.Error
}

func UpdateUser(username string, request *client.UserRequest) (int64, error) {
	user, err := FindUserByName(username)
	if err != nil {
		return 0, err
	}

	user.Email = request.Email
	user.Username = request.Username
	user.Password = request.Password

	db := database.OpenConn()
	result := db.Save(user)
	return result.RowsAffected, result.Error
}

func DeleteUserByName(username string) (int64, error) {
	db := database.OpenConn()
	result := db.Unscoped().Where("users.username = ?", username).Delete(&models.User{})
	return result.RowsAffected, result.Error
}

func SaveContact(user *models.User, contact *models.User) (int64, error) {
	db := database.OpenConn()
	user.Contacts = append(user.Contacts, contact)
	result := db.Save(user)
	return result.RowsAffected, result.Error
}
