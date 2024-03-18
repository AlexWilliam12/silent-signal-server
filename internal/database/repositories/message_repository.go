package repositories

import (
	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
)

func SavePrivateMessage(message *models.PrivateMessage) (int64, error) {
	db := database.OpenConn()
	result := db.Create(message)
	return result.RowsAffected, result.Error
}

func FetchPrivateMessages() ([]models.PrivateMessage, error) {
	db := database.OpenConn()
	var messages []models.PrivateMessage
	result := db.Find(&messages)
	return messages, result.Error
}

func FetchPendingPrivateMessages(username string) ([]models.PrivateMessage, error) {
	user, err := FindUserByName(username)
	if err != nil {
		return nil, err
	}
	db := database.OpenConn()
	var messages []models.PrivateMessage
	result := db.Where("receiver_id = ? AND is_pending = ?", user.ID, true).Find(&messages)
	return messages, result.Error
}

func UpdatePendingSituation(ids []uint) (int64, error) {
	db := database.OpenConn()
	result := db.Model(&models.PrivateMessage{}).Where("id IN ?", ids).Update("is_pending", false)
	return result.RowsAffected, result.Error
}

func SaveGroupMessage(message *models.GroupMessage) (*models.GroupMessage, error) {
	db := database.OpenConn()
	result := db.Create(message)
	return message, result.Error
}

func FetchGroupMessages() ([]models.GroupMessage, error) {
	db := database.OpenConn()
	var messages []models.GroupMessage
	result := db.Find(&messages)
	return messages, result.Error
}

func FetchPendingGroupMessages(username string) ([]models.GroupMessage, error) {
	user, err := FindUserByName(username)
	if err != nil {
		return nil, err
	}
	db := database.OpenConn()
	var messages []models.GroupMessage
	result := db.Where("id NOT IN (?)", db.Table("group_message_seen_by").Select("group_message_id").Where("user_id = ?", user.ID)).Find(&messages)
	return messages, result.Error
}
