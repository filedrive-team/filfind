package models

import "github.com/filedrive-team/filfind/backend/types"

func init() {
	autoMigrateModels = append(autoMigrateModels, &Message{})
}

type Message struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt types.UnixTime `json:"created_at" gorm:"index"`
	UpdatedAt types.UnixTime `json:"updated_at"`
	Sender    string         `json:"sender" gorm:"size:64;index:idx_msg_s_r,priority:1"`
	Recipient string         `json:"recipient" gorm:"size:64;index;index:idx_msg_s_r,priority:2"`
	Type      int8           `json:"type" gorm:"default:0;not null"` // 0:text message
	Content   string         `json:"content" gorm:"size:1024"`
	Checked   bool           `json:"checked" gorm:"default:0;not null"` // whether the Recipient has checked the message
}
