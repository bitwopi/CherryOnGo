package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager interface {
	CreateOrder(customerUUID string, status OrderStatus, shopCard *ShopCard, price float32) (*Order, error)
	GetOrder(uuid string) (*Order, error)
	UpdateOrderStatus(uuid string, status OrderStatus) (*Order, error)
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
	err := m.db.AutoMigrate(&Order{})
	if err != nil {
		log.Fatal(err)
	}
	err = m.db.AutoMigrate(&OutboxEvent{})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *PgManager) CreateOrder(customerUUID string, status OrderStatus, shopCard *ShopCard, price float32) (*Order, error) {
	order := Order{
		UUID:         uuid.NewString(),
		CustomerUUID: customerUUID,
		Status:       status,
		ShopCard:     shopCard,
		Price:        price,
	}
	if err := m.db.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *PgManager) GetOrder(uuid string) (*Order, error) {
	var order Order
	if err := m.db.Where("uuid = ?", uuid).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *PgManager) UpdateOrderStatus(uuid string, status OrderStatus) (*Order, error) {
	var order Order
	if err := m.db.
		Model(&order).
		Where("uuid = ?", uuid).
		Update("status", status).
		First(&order).Error; err != nil {
		return nil, fmt.Errorf("failed to update order: %v", err)
	}
	return &order, nil
}
