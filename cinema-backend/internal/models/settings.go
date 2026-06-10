package models

import "time"

type SystemSetting struct {
	Key       string    `json:"key" gorm:"primaryKey;not null"`
	Value     string    `json:"value" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
