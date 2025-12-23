package db

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManager interface {
	CreateUser(email string, password string, tgID int64, referral string) error
	GetUserByEmail(email string) (*User, error)
	GetUserByUUID(uuid string) (*User, error)
	GetUserByTgID(tgID int64) (*User, error)
	CheckPassword(email string, password string) error
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

func (m *PgManager) CreateUser(email string, password string, tgID int64, referral string) error {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed on create password hash: %v", err)
	}
	user := User{
		UUID:     uuid.NewString(),
		Email:    email,
		Password: string(passwordHash),
	}
	if tgID > 0 {
		user.TgID = tgID
	}
	if len(referral) == 36 {
		user.ReferralUUID = referral
	}
	err = m.db.Create(&user).Error
	if err != nil {
		return fmt.Errorf("can't create user: %v", err)
	}

	return nil
}

func (m *PgManager) CheckPassword(email string, password string) error {
	var user User
	err := m.db.
		Select("password").
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	return err
}

func (m *PgManager) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := m.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *PgManager) GetUserByUUID(uuid string) (*User, error) {
	var user User
	if err := m.db.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *PgManager) GetUserByTgID(tgID string) (*User, error) {
	var user User
	if err := m.db.Where("tg_id = ?", tgID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
