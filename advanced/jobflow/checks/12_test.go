package coursecheck_test

import (
	"jobflow/internal/protocol"
	"testing"
)

func Test12FixedHeader(t *testing.T) {
	h := protocol.Header{Lines: 3, NRefs: 1, NBF: 0x01020304, EXP: 0xffffffff}
	h.Ctx[0] = 7
	h.Dst[0] = 8
	h.Key[0] = 9
	b, err := h.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	if len(b) != 128 || b[8] != 4 || b[32] != 7 || b[64] != 8 || b[96] != 9 {
		t.Fatal("fixed offsets changed")
	}
	got, err := protocol.ParseHeader(b)
	if err != nil || got != h {
		t.Fatal(got, err)
	}
	for n := 0; n < 128; n++ {
		if _, err = protocol.ParseHeader(b[:n]); err == nil {
			t.Fatal("short header", n)
		}
	}
}
