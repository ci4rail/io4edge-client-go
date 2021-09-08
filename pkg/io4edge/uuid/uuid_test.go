package uuid

import (
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBytesToUInt64(t *testing.T) {
	b := []byte{0x12, 0x34, 0x56, 0x78, 0xab, 0xcd, 0xef, 0x11}

	v := bytesToUInt64(b)
	fmt.Printf("got %08x\n", v)
	assert.Equal(t, uint64(0x12345678abcdef11), v)
}

func TestBytesFromUInt64(t *testing.T) {
	bexp := []byte{0x12, 0x34, 0x56, 0x78, 0xab, 0xcd, 0xef, 0x11}

	b := bytesFromUint64(uint64(0x12345678abcdef11))
	assert.Equal(t, bexp, b)
}

func TestToSerial(t *testing.T) {
	u, err := uuid.NewV4()
	assert.Nil(t, err)
	b := make([]byte, 16)
	b[0] = 0x12
	b[8] = 0x34

	err = u.UnmarshalBinary(b)
	assert.Nil(t, err)

	hi, lo := ToSerial(u)
	assert.Equal(t, uint64(0x1200000000000000), hi)
	assert.Equal(t, uint64(0x3400000000000000), lo)
}

func TestFromSerial(t *testing.T) {
	u, err := FromSerial(uint64(0x1200000000000000), uint64(0x3400000000000000))
	assert.Nil(t, err)

	fmt.Printf("%s\n", u)
	assert.Equal(t, "12000000-0000-0000-3400-000000000000", u.String())
}
