package bus

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"jobflow/internal/postgres"
	"jobflow/internal/protocol"
	"log/slog"
	"time"
)

func PublishOnce(ctx context.Context, s *postgres.Store, w *kafka.Writer) error {
	events, err := s.Events(ctx, 0, true)
	if err != nil {
		return err
	}
	for _, e := range events {
		body, err := s.Payload(ctx, e.Task.Key)
		if err != nil {
			return err
		}
		h, err := protocol.ParseHeader(e.Prefix)
		if err != nil {
			return err
		}
		b := append([]byte(nil), e.Prefix...)
		b = append(b, body...)
		if int64(len(b)) > h.Size() {
			return protocol.ErrFrame
		}
		b = append(b, make([]byte, int(h.Size())-len(b))...)
		meta := e
		meta.Prefix = nil
		data, err := json.Marshal(meta)
		if err != nil {
			return err
		}
		if err = w.WriteMessages(ctx, kafka.Message{Key: []byte(e.Task.Workflow), Value: b, Headers: []kafka.Header{{Key: "event", Value: data}}}); err != nil {
			return err
		}
		if err = s.Delivered(ctx, e.ID); err != nil {
			return err
		}
		slog.Info("published", "workflow", e.Task.Workflow, "task", e.Task.ID, "event", e.ID)
	}
	return nil
}
func Relay(ctx context.Context, s *postgres.Store, brokers []string, topic string) {
	w := &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: topic, Balancer: &kafka.Hash{}, RequiredAcks: kafka.RequireAll, BatchTimeout: 10 * time.Millisecond, WriteTimeout: 2 * time.Second, MaxAttempts: 1}
	defer w.Close()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		if err := PublishOnce(ctx, s, w); err != nil && ctx.Err() == nil {
			slog.Warn("outbox delivery deferred", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
