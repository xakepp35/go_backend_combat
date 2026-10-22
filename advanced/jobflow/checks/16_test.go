package coursecheck_test

import (
	"context"
	"github.com/segmentio/kafka-go"
	"jobflow/internal/bus"
	"testing"
	"time"
)

func Test16OutboxDelivery(t *testing.T) {
	s, _ := database(t)
	submit(t, s, []byte("event"))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	bad := &kafka.Writer{Addr: kafka.TCP("127.0.0.1:1"), Topic: "none", MaxAttempts: 1, WriteTimeout: time.Millisecond}
	if err := bus.PublishOnce(ctx, s, bad); err == nil {
		t.Fatal("transport failure hidden")
	}
	_ = bad.Close()
	es, err := s.Events(ctx, 0, true)
	if err != nil || len(es) != 1 {
		t.Fatal("intent lost", len(es), err)
	}
	brokers, topic := kafkaTopic(t)
	w := &kafka.Writer{Addr: kafka.TCP(brokers...), Topic: topic, RequiredAcks: kafka.RequireAll, MaxAttempts: 1}
	defer w.Close()
	if err = bus.PublishOnce(ctx, s, w); err != nil {
		t.Fatal(err)
	}
	es, err = s.Events(ctx, 0, true)
	if err != nil || len(es) != 0 {
		t.Fatal("not confirmed", err)
	}
	r := kafka.NewReader(kafka.ReaderConfig{Brokers: brokers, Topic: topic, Partition: 0, MinBytes: 1, MaxBytes: 2 << 20})
	defer r.Close()
	m, err := r.ReadMessage(ctx)
	if err != nil || len(m.Value) < 128 {
		t.Fatal("no binary event", err)
	}
}
