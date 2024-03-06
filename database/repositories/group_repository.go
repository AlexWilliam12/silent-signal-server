package repositories

import (
	"github.com/AlexWilliam12/silent-signal/client"
	"github.com/AlexWilliam12/silent-signal/database"
	"github.com/AlexWilliam12/silent-signal/database/models"
)

// Create a group on database
func CreateGroup(group client.GroupRequest) (int64, error) {
	db := database.OpenConn()
	result := db.Select("name").Create(&models.Group{Name: group.Name})
	return result.RowsAffected, result.Error
}

// Find a group on database querying by group name
func FindGroupByName(groupName string) (*models.Group, error) {
	db := database.OpenConn()
	var group models.Group
	result := db.Where("groups.name = ?", groupName).First(&group)
	return &group, result.Error
}

// Update group informations
func UpdateGroup(group client.GroupRequest) (int64, error) {
	db := database.OpenConn()
	result := db.Save(&group)
	return result.RowsAffected, result.Error
}

// Delete a group on database where group name matches
func DeleteGroupByName(groupName string) (int64, error) {
	db := database.OpenConn()
	result := db.Unscoped().Where("groups.name = ?", groupName).Delete(&models.Group{})
	return result.RowsAffected, result.Error
}
