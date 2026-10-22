package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"jobflow/internal/postgres"
	"jobflow/internal/task"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	Base string
	HTTP *http.Client
}
type HTTPError struct {
	Status  int
	Message string
}

func (e *HTTPError) Error() string { return fmt.Sprintf("HTTP %d: %s", e.Status, e.Message) }
func New(base string) *Client {
	return &Client{Base: base, HTTP: &http.Client{Timeout: 3 * time.Second}}
}
func (c *Client) Do(ctx context.Context, method, path string, in, out any) error {
	var body []byte
	var err error
	if in != nil {
		body, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, c.Base+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return err
	}
	if resp.StatusCode >= 300 {
		return &HTTPError{resp.StatusCode, string(data)}
	}
	if out != nil {
		return json.Unmarshal(data, out)
	}
	return nil
}
func (c *Client) Claim(ctx context.Context, id string) (task.Task, error) {
	var t task.Task
	err := c.Do(ctx, "POST", "/tasks/"+url.PathEscape(id)+"/claim", nil, &t)
	return t, err
}
func (c *Client) Complete(ctx context.Context, id string, r postgres.Result) error {
	return c.Do(ctx, "POST", "/tasks/"+url.PathEscape(id)+"/result", r, nil)
}
func (c *Client) Effect(ctx context.Context, t task.Task, slot string, delta int64) (int64, error) {
	var out struct {
		Value int64 `json:"value"`
	}
	err := c.Do(ctx, "POST", "/tasks/"+url.PathEscape(t.ID)+"/effects", map[string]any{"token": t.Token, "slot": slot, "delta": delta}, &out)
	return out.Value, err
}
func (c *Client) Events(ctx context.Context, after int64) ([]postgres.Event, error) {
	var es []postgres.Event
	err := c.Do(ctx, "GET", "/internal/events?after="+strconv.FormatInt(after, 10), nil, &es)
	return es, err
}
func (c *Client) Jobs(ctx context.Context) ([]task.Task, error) {
	var ts []task.Task
	err := c.Do(ctx, "GET", "/internal/tasks", nil, &ts)
	return ts, err
}
func (c *Client) Payload(ctx context.Context, key string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.Base+"/payloads/"+url.PathEscape(key), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, &HTTPError{resp.StatusCode, "payload"}
	}
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}
