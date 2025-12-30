package db

import "time"

type ShopCard struct {
	UUID        string    `gorm:"type:varchar(36);primaryKey"`
	Name        string    `gorm:"type:varchar(100);not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	Description string    `gorm:"type:varchar(500)"`
	Category    string    `gorm:"type:varchar(100)"`
	Price       *float64  `gorm:"type:decimal(10,2);"`
	Visible     bool      `gorm:"default:true"`
	CoverURL    *string   `gorm:"type:varchar(100)"`
}

var (
	Categories = []string{
		"sert",
		"delivery",
		"disk",
		"subscription",
		"other",
	}
)
