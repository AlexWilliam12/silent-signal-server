package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email            string `gorm:"uniqueIndex;not null;index"`
	Username         string `gorm:"uniqueIndex;not null;index"`
	Password         string `gorm:"not null"`
	Picture          string
	SentMessages     []PrivateMessage `gorm:"foreignKey:SenderID"`
	ReceivedMessages []PrivateMessage `gorm:"foreignKey:ReceiverID"`
	Groups           []*Group         `gorm:"many2many:group_members;"`
	Contacts         []*User          `gorm:"many2many:user_contacts;association_jointable_foreignkey:contact_id"`
}

type PrivateMessage struct {
	gorm.Model
	SenderID   uint
	ReceiverID uint
	Sender     User   `gorm:"foreignKey:SenderID"`
	Receiver   User   `gorm:"foreignKey:ReceiverID"`
	Type       string `gorm:"not null"`
	Content    string `gorm:"not null"`
	IsPending  bool   `gorm:"not null"`
}

type Group struct {
	gorm.Model
	Name          string `gorm:"not null;uniqueIndex;index"`
	Description   string
	Picture       string
	CreatorID     uint
	Creator       User    `gorm:"not null;foreignKey:CreatorID"`
	Members       []*User `gorm:"many2many:group_members;"`
	GroupMessages []GroupMessage
}

type GroupMessage struct {
	gorm.Model
	SenderID uint
	GroupID  uint
	Group    Group   `gorm:"foreignKey:GroupID"`
	Sender   User    `gorm:"foreignKey:SenderID"`
	SeenBy   []*User `gorm:"many2many:group_message_seen_by;"`
	Type     string  `gorm:"not null"`
	Content  string  `gorm:"not null"`
}
