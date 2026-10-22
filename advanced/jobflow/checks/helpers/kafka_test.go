package coursecheck_test

import (
	"context"
	"github.com/segmentio/kafka-go"
	"jobflow/internal/task"
	"os"
	"strings"
	"testing"
	"time"
)

func kafkaTopic(t *testing.T) ([]string, string) {
	t.Helper()
	brokers := strings.Split(os.Getenv("JF_TEST_KAFKA_BROKERS"), ",")
	if brokers[0] == "" {
		t.Fatal("Kafka infrastructure required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	name := "check-" + task.NewID()
	if err = conn.CreateTopics(kafka.TopicConfig{Topic: name, NumPartitions: 1, ReplicationFactor: 1}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		c, e := kafka.Dial("tcp", brokers[0])
		if e == nil {
			_ = c.DeleteTopics(name)
			_ = c.Close()
		}
	})
	return brokers, name
}
