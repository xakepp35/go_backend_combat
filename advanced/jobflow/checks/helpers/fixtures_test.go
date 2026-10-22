package coursecheck_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"jobflow/internal/protocol"
	"os"
	"testing"
)

func Test13PublishedFixtures(t *testing.T) {
	data, err := os.ReadFile("fixtures/frames.json")
	if err != nil {
		t.Fatal(err)
	}
	var vectors []struct {
		Name, Hex string
		Valid     bool
		Body      string `json:"body_hex"`
		Refs      []uint64
		Deps      []string
	}
	if err = json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	for _, v := range vectors {
		t.Run(v.Name, func(t *testing.T) {
			raw, err := hex.DecodeString(v.Hex)
			if err != nil {
				t.Fatal(err)
			}
			f, err := protocol.Unmarshal(raw)
			if !v.Valid {
				if err == nil {
					t.Fatal("malformed fixture accepted")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			body, err := hex.DecodeString(v.Body)
			if err != nil || !bytes.Equal(body, f.Body) {
				t.Fatal("body changed", err)
			}
			if len(v.Refs) != len(f.Refs) || len(v.Deps) != len(f.Deps) {
				t.Fatal("links changed")
			}
			for i, r := range v.Refs {
				if f.Refs[i] != r {
					t.Fatal("ref byte order")
				}
			}
			for i, d := range v.Deps {
				if hex.EncodeToString(f.Deps[i][:]) != d {
					t.Fatal("dep byte order")
				}
			}
			canonical, err := protocol.Marshal(f)
			if err != nil || !bytes.Equal(raw, canonical) {
				t.Fatal("fixture roundtrip", err)
			}
		})
	}
}
