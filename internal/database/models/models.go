package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username         string `gorm:"uniqueIndex;not null;index"`
	Password         string `gorm:"not null"`
	Picture          string
	SentMessages     []PrivateMessage `gorm:"foreignKey:SenderID"`
	ReceivedMessages []PrivateMessage `gorm:"foreignKey:ReceiverID"`
	Groups           []Group          `gorm:"many2many:user_groups;"`
}

type PrivateMessage struct {
	gorm.Model
	SenderID   uint
	ReceiverID uint
	Sender     User   `gorm:"foreignKey:SenderID"`
	Receiver   User   `gorm:"foreignKey:ReceiverID"`
	Data       string `gorm:"not null"`
	IsPending  bool   `gorm:"not null"`
}

type Group struct {
	gorm.Model
	Name          string `gorm:"not null;uniqueIndex;index"`
	Description   string
	Picture       string
	CreatorID     uint
	Creator       User   `gorm:"not null;foreignKey:CreatorID"`
	Members       []User `gorm:"many2many:user_groups;"`
	GroupMessages []GroupMessage
}

type GroupMessage struct {
	gorm.Model
	SenderID uint
	GroupID  uint   `gorm:"foreignKey:GroupID"`
	Sender   User   `gorm:"foreignKey:SenderID"`
	Data     string `gorm:"not null"`
}