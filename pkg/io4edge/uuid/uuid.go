package uuid

import (
	"github.com/gofrs/uuid"
)

// ToSerial converts the UUID u into hi and lo
func ToSerial(u uuid.UUID) (hi uint64, lo uint64) {

	hi = bytesToUInt64(u.Bytes()[0:8])
	lo = bytesToUInt64(u.Bytes()[8:16])
	return hi, lo
}

func bytesToUInt64(b []byte) uint64 {
	var ret = uint64(0)

	for i := 0; i < 8; i++ {
		ret <<= 8
		ret += uint64(b[i])
	}
	return ret
}

func bytesFromUint64(v uint64) []byte {
	b := make([]byte, 8)

	for i := 0; i < 8; i++ {
		b[i] = byte(v >> 56)
		v <<= 8
	}
	return b
}

// FromSerial hi and lo into an UUID
func FromSerial(hi uint64, lo uint64) (uuid.UUID, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return uuid.UUID{}, err
	}
	b := make([]byte, 16)

	copy(b[0:8], bytesFromUint64(hi))
	copy(b[8:16], bytesFromUint64(lo))

	err = u.UnmarshalBinary(b)
	if err != nil {
		return uuid.UUID{}, err
	}

	return u, nil
}
