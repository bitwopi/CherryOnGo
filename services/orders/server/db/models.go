package db

import "time"

type Order struct {
	UUID         string      `gorm:"type:varchar(36);primaryKey"`
	CreatedAt    time.Time   `gorm:"autoCreateTime"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime"`
	CustomerUUID string      `gorm:"type:varchar(36)"`
	Status       OrderStatus `gorm:"type:varchar(20);default:'NEW';check:status IN ('NEW','UNPAID','PAID','CANCELLED','REFUNDED')"`
	ShopCard     ShopCard    `gorm:"type:jsonb;serializer:json" json:"shop_card"`
	Price        float32     `gorm:"type:decimal(10,2)"`
}

type OrderStatus string

const (
	StatusNew       OrderStatus = "NEW"
	StatusUnpaid    OrderStatus = "UNPAID"
	StatusPaid      OrderStatus = "PAID"
	StatusCancelled OrderStatus = "CANCELLED"
	StatusRefunded  OrderStatus = "REFUNDED"
)

type ShopCard struct {
	UUID        string
	Name        string
	CreatedAt   time.Time
	Description string
	Category    string
	Price       *float32
	Visible     bool
	CoverURL    *string
}

type OutboxEvent struct {
	UUID          string     `gorm:"varchar(36);primaryKey"`
	AggregateID   string     `gorm:"varchar(36);not null"`
	AggregateType string     `gorm:"text;not null"`
	EventType     string     `gorm:"text;not null"`
	Payload       []byte     `gorm:"bytea;not null"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	ProcessedAt   *time.Time `gorm:"timestamp"`
}
