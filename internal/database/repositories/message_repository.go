package repositories

import (
	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
)

func SavePrivateMessage(message *models.PrivateMessage) (*models.PrivateMessage, error) {
	db := database.OpenConn()
	result := db.Create(message)
	return message, result.Error
}

func FetchPendingMessages() ([]models.PrivateMessage, error) {
	db := database.OpenConn()
	var messages []models.PrivateMessage
	result := db.Where("is_pending = TRUE").Find(&messages)
	return messages, result.Error
}

func UpdatePendingSituation(message *models.PrivateMessage) (int64, error) {
	db := database.OpenConn()
	result := db.Save(&message)
	return result.RowsAffected, result.Error
}
