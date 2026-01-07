package outbox

import (
	"context"
	"orders/server/db"
	"time"

	"github.com/segmentio/kafka-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OutboxRepository interface {
	GetUnprocessed(ctx context.Context, limit int) ([]db.OutboxEvent, error)
	MarkProcessed(ctx context.Context, uuid string) error
}

type Publisher interface {
	Publish(ctx context.Context, topic string, key []byte, value []byte) error
}

type KafkaPublisher struct {
	Publisher
	Addr string
}

func NewKafkaPublisher(addr string) *KafkaPublisher {
	return &KafkaPublisher{Addr: addr}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic string, key []byte, value []byte) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		writer := kafka.Writer{
			Addr:                   kafka.TCP(p.Addr),
			Topic:                  topic,
			Balancer:               &kafka.RoundRobin{},
			AllowAutoTopicCreation: true,
		}

		err := writer.WriteMessages(ctx, kafka.Message{
			Key:   key,
			Topic: topic,
			Value: value,
		})
		if err != nil {
			return err
		}
		return writer.Close()
	}

}

type PgRepository struct {
	db *gorm.DB
}

func NewPgRepository(dsn string) (*PgRepository, error) {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PgRepository{db: conn}, nil
}

func (p *PgRepository) GetUnprocessed(ctx context.Context, limit int) ([]db.OutboxEvent, error) {
	var events []db.OutboxEvent
	err := p.db.
		WithContext(ctx).
		Model(&db.OutboxEvent{}).
		Where("processed_at IS NULL").
		Order("created_at ASC").
		Limit(limit).
		Find(&events).
		Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (p *PgRepository) MarkProcessed(ctx context.Context, uuid string) error {
	return p.db.
		WithContext(ctx).
		Model(&db.OutboxEvent{}).
		Where("uuid = ?", uuid).
		Update("processed_at", time.Now()).
		Error
}
