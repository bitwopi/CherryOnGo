package db

import (
	"fmt"
	"log"
	"payment/server/sdk"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager interface {
	CreatePayment(payment sdk.Payment) (*Payment, error)
	CreateRefund(refund sdk.Refund) (*Refund, error)
	GetPayment(pUUID string) (*Payment, error)
	GetRefund(rUUID string) (*Refund, error)
	UpdatePaymentStatus(pUUID string, status string) (*Payment, error)
	UpdateRefundStatus(rUUID string, status string) (*Refund, error)
	Migrate()
}

type PgManager struct {
	DBManager
	db *gorm.DB
}

func NewManager(dsn string) (*PgManager, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PgManager{db: conn}, nil
}

func (m *PgManager) Migrate() {
	err := m.db.AutoMigrate(&Payment{})
	if err != nil {
		log.Fatal(err)
	}
	err = m.db.AutoMigrate(&Refund{})
	if err != nil {
		log.Fatal(err)
	}
}

func (p *PgManager) CreatePayment(payment sdk.Payment) (*Payment, error) {
	email := payment.Metadata["email"].(string)
	tgID, ok := payment.Metadata["tg_id"].(string)
	var tgIDInt int64
	if ok {
		temp, _ := strconv.Atoi(tgID)
		tgIDInt = int64(temp)
	}
	cUUID := payment.Metadata["customer_uuid"].(string)
	oUUID := payment.Metadata["oreder_uuid"].(string)
	newPayment := Payment{
		UUID:         uuid.NewString(),
		CreatedAt:    *payment.CreatedAt,
		Email:        &email,
		TgID:         &tgIDInt,
		CustomerUUID: cUUID,
		OrderUUID:    oUUID,
		Amount:       payment.Amount.Value,
		Income:       payment.Income.Value,
		Provider:     payment.Provider,
		Status:       PaymnetPending,
	}
	err := p.db.Create(&newPayment).Error
	if err != nil {
		return nil, err
	}

	return &newPayment, nil
}

func (p *PgManager) GetPayment(pUUID string) (*Payment, error) {
	var payment Payment
	if err := p.db.Where("uuid = ?", pUUID).First(&payment).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (p *PgManager) UpdatePaymentStatus(pUUID string, status PaymentStatus) (*Payment, error) {
	var payment Payment
	if err := p.db.
		Model(&payment).
		Where("uuid = ?", pUUID).
		Update("status", status).
		First(&payment).Error; err != nil {
		return nil, fmt.Errorf("failed to update payment: %v", err)
	}
	return &payment, nil
}

func (p *PgManager) CreateRefund(refund sdk.Refund) (*Refund, error) {
	cUUID := refund.Metadata["customer_uuid"].(string)
	newRefund := Refund{
		UUID:         uuid.NewString(),
		CreatedAt:    *refund.CreatedAt,
		CustomerUUID: cUUID,
		Amount:       refund.Amount.Value,
		Status:       RefundPending,
	}
	err := p.db.Create(&newRefund).Error
	if err != nil {
		return nil, err
	}

	return &newRefund, nil
}

func (p *PgManager) GetRefund(rUUID string) (*Refund, error) {
	var refund Refund
	if err := p.db.Where("uuid = ?", rUUID).First(&refund).Error; err != nil {
		return nil, err
	}
	return &refund, nil
}

func (p *PgManager) UpdateRefundStatus(rUUID string, status RefundStatus) (*Refund, error) {
	var refund Refund
	if err := p.db.
		Model(&refund).
		Where("uuid = ?", rUUID).
		Update("status", status).
		First(&refund).Error; err != nil {
		return nil, fmt.Errorf("failed to update refund: %v", err)
	}
	return &refund, nil
}
