package protocol

import (
	"bytes"
	"encoding/binary"
	"github.com/zeebo/blake3"
	"io"
)

type Frame struct {
	Header Header
	Refs   []uint64
	Deps   [][16]byte
	Body   []byte
}

func Marshal(f Frame) ([]byte, error) {
	if len(f.Refs) > 65535 || len(f.Deps) > 65535 {
		return nil, ErrFrame
	}
	n := HeaderSize + len(f.Refs)*8 + len(f.Deps)*16 + len(f.Body)
	padding := (64 - n%64) % 64
	if n > MaxFrame-padding {
		return nil, ErrFrame
	}
	f.Header.Lines = uint32((n + padding) / 64)
	f.Header.NRefs = uint16(len(f.Refs))
	f.Header.NDeps = uint16(len(f.Deps))
	f.Header.Key = blake3.Sum256(f.Body)
	b, err := f.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}
	for _, ref := range f.Refs {
		b = binary.LittleEndian.AppendUint64(b, ref)
	}
	for _, dep := range f.Deps {
		b = append(b, dep[:]...)
	}
	b = append(b, f.Body...)
	b = append(b, make([]byte, padding)...)
	return b, nil
}

// Prefix reads metadata from an already bounded whole record; it never decodes body.
func Prefix(b []byte) (Frame, error) {
	var f Frame
	h, err := ParseHeader(b)
	if err != nil {
		return f, err
	}
	if int64(len(b)) != h.Size() {
		return f, ErrFrame
	}
	f.Header = h
	off := HeaderSize
	for range int(h.NRefs) {
		f.Refs = append(f.Refs, binary.LittleEndian.Uint64(b[off:off+8]))
		off += 8
	}
	for range int(h.NDeps) {
		var d [16]byte
		copy(d[:], b[off:off+16])
		f.Deps = append(f.Deps, d)
		off += 16
	}
	return f, nil
}
func Unmarshal(b []byte) (Frame, error) {
	f, err := Prefix(b)
	if err != nil {
		return f, err
	}
	tail := b[HeaderSize+int(f.Header.LinksSize()):]
	// Hash the common prefix once, then extend one byte at a time. Sum does not reset.
	start := max(0, len(tail)-63)
	h := blake3.New()
	_, _ = h.Write(tail[:start])
	match := -1
	for length := start; length <= len(tail); length++ {
		pad := tail[length:]
		if len(pad) < 64 && bytes.Count(pad, []byte{0}) == len(pad) && bytes.Equal(h.Sum(nil), f.Header.Key[:]) {
			if match != -1 {
				return Frame{}, ErrFrame
			}
			match = length
		}
		if length < len(tail) {
			_, _ = h.Write(tail[length : length+1])
		}
	}
	if match < 0 {
		return Frame{}, ErrFrame
	}
	f.Body = bytes.Clone(tail[:match])
	return f, nil
}
func Read(r io.Reader) (Frame, error) {
	hb := make([]byte, HeaderSize)
	_, err := io.ReadFull(r, hb)
	if err != nil {
		return Frame{}, err
	}
	h, err := ParseHeader(hb)
	if err != nil {
		return Frame{}, err
	}
	b := make([]byte, int(h.Size()))
	copy(b, hb)
	if _, err = io.ReadFull(r, b[HeaderSize:]); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return Frame{}, err
	}
	return Unmarshal(b)
}
func Skip(r io.Reader) error {
	var first [4]byte
	if _, err := io.ReadFull(r, first[:]); err != nil {
		return err
	}
	lines := binary.LittleEndian.Uint32(first[:])
	if lines < 2 || int64(lines)*64 > MaxFrame {
		return ErrFrame
	}
	_, err := io.CopyN(io.Discard, r, int64(lines)*64-4)
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}
