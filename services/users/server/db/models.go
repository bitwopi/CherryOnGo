package db

import (
	"time"
)

type User struct {
	UUID     string `gorm:"type:varchar(36);primaryKey"`
	Email    string `gorm:"type:varchar(256);uniqueIndex;not null"`
	TgID     int64  `gorm:"uniqueIndex"`
	Password string `gorm:"type:varchar(256);not null"`

	CreatedAt time.Time `gorm:"autoCreateTime"`

	Active       bool   `gorm:"default:true"`
	ReferralUUID string `gorm:"type:varchar(36)"`
	Trial        bool   `gorm:"default:false"`
}
