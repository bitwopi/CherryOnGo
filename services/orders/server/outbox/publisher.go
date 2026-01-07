package outbox

import (
	"context"
	"errors"
	"time"
)

type OutboxPublisher struct {
	repo      OutboxRepository
	publisher Publisher
	interval  time.Duration
}

func NewOutboxPublisher(dsn string, brockerAddr string, tickInterval time.Duration) *OutboxPublisher {
	repository, err := NewPgRepository(dsn)
	if err != nil {
		panic(err)
	}
	pub := NewKafkaPublisher(brockerAddr)
	return &OutboxPublisher{
		repo:      repository,
		publisher: pub,
		interval:  tickInterval,
	}
}

func (p *OutboxPublisher) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.publishBatch(ctx)
		}
	}
}

func (p *OutboxPublisher) publishBatch(ctx context.Context) {
	events, err := p.repo.GetUnprocessed(ctx, 10)
	if err != nil {
		return
	}

	for _, e := range events {
		err := p.publisher.Publish(ctx, e.AggregateType+"-service", []byte(e.AggregateID), e.Payload)
		if errors.Is(err, ctx.Err()) {
			return
		} else if err != nil {
			continue
		}

		_ = p.repo.MarkProcessed(ctx, e.UUID)
	}
}
