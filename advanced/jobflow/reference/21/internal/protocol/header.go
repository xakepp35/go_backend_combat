// Package protocol implements the fixed 128-byte framing contract, without unsafe.
package protocol

import (
	"encoding/binary"
	"errors"
)

const HeaderSize = 128
const MaxFrame = 1 << 20

var ErrFrame = errors.New("invalid frame")

type Header struct {
	Lines                    uint32
	NRefs, NDeps             uint16
	NBF, EXP, TX, RX, FX, SX uint32
	Ctx, Dst, Key            [32]byte
}

func (h Header) Size() int64      { return int64(h.Lines) * 64 }
func (h Header) LinksSize() int64 { return int64(h.NRefs)*8 + int64(h.NDeps)*16 }
func (h Header) Validate() error {
	if h.Lines < 2 || h.Size() > MaxFrame || HeaderSize+h.LinksSize() > h.Size() {
		return ErrFrame
	}
	if h.EXP != 0 && h.NBF >= h.EXP {
		return ErrFrame
	}
	return nil
}
func ParseHeader(b []byte) (Header, error) {
	var h Header
	if len(b) < HeaderSize {
		return h, ErrFrame
	}
	h.Lines = binary.LittleEndian.Uint32(b[0:4])
	h.NRefs = binary.LittleEndian.Uint16(b[4:6])
	h.NDeps = binary.LittleEndian.Uint16(b[6:8])
	h.NBF = binary.LittleEndian.Uint32(b[8:12])
	h.EXP = binary.LittleEndian.Uint32(b[12:16])
	h.TX = binary.LittleEndian.Uint32(b[16:20])
	h.RX = binary.LittleEndian.Uint32(b[20:24])
	h.FX = binary.LittleEndian.Uint32(b[24:28])
	h.SX = binary.LittleEndian.Uint32(b[28:32])
	copy(h.Ctx[:], b[32:64])
	copy(h.Dst[:], b[64:96])
	copy(h.Key[:], b[96:128])
	return h, h.Validate()
}
func (h Header) MarshalBinary() ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}
	b := make([]byte, HeaderSize)
	binary.LittleEndian.PutUint32(b[0:4], h.Lines)
	binary.LittleEndian.PutUint16(b[4:6], h.NRefs)
	binary.LittleEndian.PutUint16(b[6:8], h.NDeps)
	for i, v := range []uint32{h.NBF, h.EXP, h.TX, h.RX, h.FX, h.SX} {
		binary.LittleEndian.PutUint32(b[8+i*4:12+i*4], v)
	}
	copy(b[32:64], h.Ctx[:])
	copy(b[64:96], h.Dst[:])
	copy(b[96:128], h.Key[:])
	return b, nil
}
