package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"testing"
)

func TestHeaderOffsets(t *testing.T) {
	h := Header{Lines: 4, NRefs: 1, NDeps: 1, NBF: 0x01020304, EXP: 0xffffffff, TX: 5, RX: 6, FX: 7, SX: 8}
	h.Ctx[0] = 9
	h.Dst[0] = 10
	h.Key[0] = 11
	b, err := h.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 128 || b[8] != 4 || b[9] != 3 || b[32] != 9 || b[64] != 10 || b[96] != 11 {
		t.Fatal("wire offsets changed")
	}
	got, err := ParseHeader(b)
	if err != nil || got != h {
		t.Fatalf("roundtrip: %#v %v", got, err)
	}
}
func TestHeaderBounds(t *testing.T) {
	for n := 0; n < 128; n++ {
		if _, err := ParseHeader(make([]byte, n)); err == nil {
			t.Fatalf("accepted %d bytes", n)
		}
	}
	for _, h := range []Header{{Lines: 0}, {Lines: 1}, {Lines: 0xffffffff}, {Lines: 2, NRefs: 1}, {Lines: 2, NDeps: 1}, {Lines: 2, NBF: 10, EXP: 10}} {
		if _, err := h.MarshalBinary(); err == nil {
			t.Fatalf("accepted %#v", h)
		}
	}
}
func TestFrames(t *testing.T) {
	for _, body := range [][]byte{nil, {}, []byte("a"), {0}, {1, 0, 0}, bytes.Repeat([]byte{0}, 64), bytes.Repeat([]byte{2}, 65), bytes.Repeat([]byte{3}, 1024)} {
		f := Frame{Body: body, Refs: []uint64{42}, Deps: [][16]byte{{1}}}
		b, err := Marshal(f)
		if err != nil {
			t.Fatal(err)
		}
		got, err := Unmarshal(b)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(body, got.Body) || len(b)%64 != 0 || got.Refs[0] != 42 || got.Deps[0][0] != 1 {
			t.Fatal("frame changed")
		}
		prefix, err := Prefix(b)
		if err != nil || len(prefix.Body) != 0 {
			t.Fatal("prefix decoded body", err)
		}
		for n := 0; n < len(b); n++ {
			if _, err := Read(bytes.NewReader(b[:n])); err == nil {
				t.Fatalf("accepted truncation %d", n)
			}
		}
		corrupt := bytes.Clone(b)
		corrupt[96] ^= 1
		if _, err := Unmarshal(corrupt); err == nil {
			t.Fatal("accepted corrupt hash")
		}
	}
}
func TestEmptyGolden(t *testing.T) {
	b, err := Marshal(Frame{})
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 128 || b[0] != 2 || hex.EncodeToString(b[96:]) != "af1349b9f5f9a1a6a0404dea36dcc9499bcb25c9adc112b7cc9a93cae41f3262" {
		t.Fatal("BLAKE3/empty golden changed")
	}
}

type byteReader struct{ io.Reader }

func (r byteReader) Read(p []byte) (int, error) {
	if len(p) > 1 {
		p = p[:1]
	}
	return r.Reader.Read(p)
}
func TestStream(t *testing.T) {
	a, _ := Marshal(Frame{Body: []byte("first")})
	b, _ := Marshal(Frame{Body: []byte{4, 0}})
	r := byteReader{bytes.NewReader(append(a, b...))}
	if err := Skip(r); err != nil {
		t.Fatal(err)
	}
	f, err := Read(r)
	if err != nil || !bytes.Equal(f.Body, []byte{4, 0}) {
		t.Fatal(err)
	}
	if _, err = Read(r); !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}
	if err = Skip(bytes.NewReader(a[:4])); !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Fatal("skip truncated", err)
	}
	if err = Skip(bytes.NewReader([]byte{1, 0, 0, 0})); err == nil {
		t.Fatal("skip short lines")
	}
}
func TestLimitsAndCorruption(t *testing.T) {
	if _, err := Marshal(Frame{Body: make([]byte, MaxFrame)}); err == nil {
		t.Fatal("unbounded body")
	}
	if _, err := Marshal(Frame{Refs: make([]uint64, 65536)}); err == nil {
		t.Fatal("count overflow")
	}
	if _, err := Marshal(Frame{Deps: make([][16]byte, 65536)}); err == nil {
		t.Fatal("count overflow")
	}
	b, _ := Marshal(Frame{Body: []byte("x")})
	binary.LittleEndian.PutUint16(b[4:6], 65535)
	if _, err := Unmarshal(b); err == nil {
		t.Fatal("links outside frame")
	}
	b, _ = Marshal(Frame{Body: []byte("x")})
	b[len(b)-1] = 1
	if _, err := Unmarshal(b); err == nil {
		t.Fatal("nonzero padding")
	}
	b, _ = Marshal(Frame{})
	if _, err := Unmarshal(append(b, 0)); err == nil {
		t.Fatal("trailing bytes")
	}
	b, _ = Marshal(Frame{})
	binary.LittleEndian.PutUint32(b[:4], 0xffffffff)
	if _, err := Read(bytes.NewReader(b)); err == nil {
		t.Fatal("oversized allocation")
	}
}
func FuzzFrame(f *testing.F) {
	b, _ := Marshal(Frame{Body: []byte{1, 0, 0}})
	f.Add(b)
	f.Add([]byte{})
	f.Add(make([]byte, 128))
	f.Fuzz(func(t *testing.T, b []byte) {
		if len(b) > MaxFrame+1 {
			return
		}
		frame, err := Unmarshal(b)
		if err == nil {
			out, err := Marshal(frame)
			if err != nil || !bytes.Equal(out, b) {
				t.Fatal("noncanonical roundtrip", err)
			}
		}
	})
}
