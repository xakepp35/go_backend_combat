package coursecheck_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"jobflow/internal/postgres"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func database(t *testing.T) (*postgres.Store, string) {
	t.Helper()
	dsn := os.Getenv("JF_TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("JF_TEST_DATABASE_URL required: run coursecheck, infrastructure cannot be skipped")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	admin, err := postgres.Open(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	schema := "check_" + strings.ToLower(rand.Text())
	if _, err = admin.Pool.Exec(ctx, fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)); err != nil {
		admin.Close()
		t.Fatal(err)
	}
	u, err := url.Parse(dsn)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()
	isolated := u.String()
	s, err := postgres.Open(ctx, isolated)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		s.Close()
		_, _ = admin.Pool.Exec(context.Background(), fmt.Sprintf(`DROP SCHEMA "%s" CASCADE`, schema))
		admin.Close()
	})
	if err = s.Migrate(ctx, "../migrations"); err != nil {
		t.Fatal(err)
	}
	return s, isolated
}
