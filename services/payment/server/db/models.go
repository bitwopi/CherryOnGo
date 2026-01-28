package db

import "time"

type Payment struct {
	UUID         string        `gorm:"type:varchar(36);primaryKey"`
	CreatedAt    time.Time     `gorm:"type:timestamptz"`
	Email        *string       `gorm:"type:varchar(256)"`
	TgID         *int64        `gorm:"type:int"`
	CustomerUUID string        `gorm:"type:varchar(36)"`
	OrderUUID    string        `gorm:"type:varchar(36)"`
	Amount       float64       `gorm:"type:decimal"`
	Income       float64       `gorm:"type:decimal"`
	Provider     string        `gorm:"type:varchar(18)"`
	Status       PaymentStatus `gorm:"type:varchar(18)"`
}

type PaymentStatus string

const (
	PaymnetPending    PaymentStatus = "PENDING"
	WaitingForCapture PaymentStatus = "WAITING_FOR_CAPTURE"
	PaymentSucceeded  PaymentStatus = "SUCCEEDED"
	PaymentCanceled   PaymentStatus = "CANCELED"
)

type Refund struct {
	UUID         string       `gorm:"type:varchar(36);primaryKey"`
	CreatedAt    time.Time    `gorm:"type:timestamptz"`
	Amount       float64      `gorm:"type:decimal"`
	PaymentUUID  string       `gorm:"type:varchar(36)"`
	CustomerUUID string       `gorm:"type:varchar(36)"`
	Status       RefundStatus `gorm:"type:varchar(10)"`
}

type RefundStatus string

const (
	RefundPending   RefundStatus  = "PENDING"
	RefundSucceeded RefundStatus  = "SUCCEEDED"
	RefundCanceled  PaymentStatus = "CANCELED"
)
