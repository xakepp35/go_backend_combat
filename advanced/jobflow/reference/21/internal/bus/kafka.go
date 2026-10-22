package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"jobflow/internal/postgres"
	"jobflow/internal/projection"
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
func ConsumeOne(ctx context.Context, r *kafka.Reader, p projection.Projection) error {
	m, err := r.FetchMessage(ctx)
	if err != nil {
		return err
	}
	f, err := protocol.Prefix(m.Value)
	if err != nil {
		return err
	}
	var e postgres.Event
	found := false
	for _, h := range m.Headers {
		if h.Key == "event" {
			if err = json.Unmarshal(h.Value, &e); err != nil {
				return err
			}
			found = true
		}
	}
	if !found || e.ID == "" || e.Task.Key != fmt.Sprintf("%x", f.Header.Key) {
		return protocol.ErrFrame
	}
	e.Prefix = m.Value[:protocol.HeaderSize+int(f.Header.LinksSize())]
	if err = p.Apply(ctx, e); err != nil {
		return err
	}
	if err = r.CommitMessages(ctx, m); err != nil {
		return err
	}
	slog.Info("projected", "workflow", e.Task.Workflow, "task", e.Task.ID, "event", e.ID)
	return nil
}
func Consume(ctx context.Context, p projection.Projection, brokers []string, topic, group string) {
	// On Redis failure recreate the reader: the failed, uncommitted record must be fetched again.
	for ctx.Err() == nil {
		r := kafka.NewReader(kafka.ReaderConfig{Brokers: brokers, Topic: topic, GroupID: group, MinBytes: 1, MaxBytes: 2 << 20, CommitInterval: 0, StartOffset: kafka.FirstOffset})
		for ctx.Err() == nil {
			if err := ConsumeOne(ctx, r, p); err != nil {
				if ctx.Err() == nil {
					slog.Warn("projection delivery deferred", "error", err)
				}
				break
			}
		}
		_ = r.Close()
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}
