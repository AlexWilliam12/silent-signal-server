package repositories

import (
	"github.com/AlexWilliam12/silent-signal/internal/client"
	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/AlexWilliam12/silent-signal/internal/database/models"
)

func CreateGroup(request *client.GroupRequest, creator *models.User) (int64, error) {
	db := database.OpenConn()
	result := db.Create(&models.Group{Name: request.Name, Description: request.Description, Creator: *creator})
	return result.RowsAffected, result.Error
}

func FindAllGroups() ([]models.Group, error) {
	db := database.OpenConn()
	var groups []models.Group
	result := db.Find(&groups)
	return groups, result.Error
}

func FindGroupByName(groupName string) (*models.Group, error) {
	db := database.OpenConn()
	var group models.Group
	result := db.Where("groups.name = ?", groupName).First(&group)
	return &group, result.Error
}

func UpdateGroup(request *client.GroupRequest, group *models.Group) (int64, error) {
	group.Name = request.Name
	group.Description = request.Description
	db := database.OpenConn()
	result := db.Save(&group)
	return result.RowsAffected, result.Error
}

func DeleteGroupByName(groupName string) (int64, error) {
	db := database.OpenConn()
	result := db.Unscoped().Where("groups.name = ?", groupName).Delete(&models.Group{})
	return result.RowsAffected, result.Error
}
