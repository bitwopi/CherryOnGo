package db

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager interface {
	CreateShopCard(name string, description string, category string, price *float64, visible bool, coverURL *string) (*ShopCard, error)
	UpdateShopCard(uuid string, name string, description string, category string, price *float64, visible bool, coverURL *string) (*ShopCard, error)
	GetShopCard(uuid string) (*ShopCard, error)
	DeleteShopCard(uuid string) error
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
	err := m.db.AutoMigrate(&ShopCard{})
	if err != nil {
		log.Fatal(err)
	}
}

func (m *PgManager) CreateShopCard(
	name string,
	description string,
	category string,
	price *float64,
	visible bool,
	coverURL *string) (*ShopCard, error) {
	card := ShopCard{
		UUID:        uuid.NewString(),
		Name:        name,
		Description: description,
		Price:       price,
		Visible:     visible,
		CoverURL:    coverURL,
	}
	if !contains(Categories, category) {
		return nil, errors.New("invalid category name")
	}
	card.Category = category
	if err := m.db.Create(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (m *PgManager) UpdateShopCard(
	uuid string,
	name string,
	description string,
	category string,
	price *float64,
	visible bool,
	coverURL *string) (*ShopCard, error) {
	card := ShopCard{
		UUID:        uuid,
		Name:        name,
		Description: description,
		Price:       price,
		Visible:     visible,
		CoverURL:    coverURL,
	}
	if !contains(Categories, category) {
		return nil, errors.New("invalid category name")
	}
	if err := m.db.Model(&card).Updates(card).Error; err != nil {
		return nil, fmt.Errorf("failed to create card: %v", err)
	}
	return &card, nil
}

func (m *PgManager) GetShopCard(uuid string) (*ShopCard, error) {
	var card ShopCard
	if err := m.db.Model(&ShopCard{}).Where("uuid = ?", uuid).First(&card).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (m *PgManager) DeleteShopCard(uuid string) error {
	return m.db.Where("uuid = ?", uuid).Delete(&ShopCard{}).Error
}

func contains[T comparable](arr []T, target T) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}
