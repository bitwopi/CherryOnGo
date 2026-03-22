package db

import (
	"time"
)

type User struct {
	UUID         string    `gorm:"type:varchar(36);primaryKey"`
	Email        *string   `gorm:"type:varchar(256);uniqueIndex"`
	TgID         *int64    `gorm:"type:int;uniqueIndex"`
	Password     *string   `gorm:"type:varchar(256)"`
	Username     *string   `gorm:"type:varchar(100)"`
	FirstName    *string   `gorm:"type:varchar(100)"`
	PhotoURL     *string   `gorm:"type:varchar(256)"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Active       bool      `gorm:"default:true"`
	ReferralUUID *string   `gorm:"type:varchar(36)"`
	Trial        bool      `gorm:"default:false"`
	Roles        *[]Role   `gorm:"type:jsonb;serializer:json" json:"roles"`
}

type Role string

const (
	UserRole Role = "USER"
	Admin    Role = "ADMIN"
	Vless    Role = "VLESS"
	Anime    Role = "ANIME"
	Tables   Role = "TABLES"
	Cloud    Role = "REFUNDED"
)

func RolesToStrings(roles []Role) []string {
	result := make([]string, len(roles))
	for i, r := range roles {
		result[i] = string(r)
	}
	return result
}
