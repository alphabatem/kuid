package kuid

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

// KUID represents a compressed universally unique identifier
type KUID struct {
	msb uint64 // most significant bits
	lsb uint64 // least significant bits
}

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base        = uint64(len(base62Chars))
	size        = 11 // size of each encoded long
)

var (
	ErrInvalidLength = errors.New("invalid KUID string length")
	ErrInvalidChar   = errors.New("invalid character in KUID string")
	ErrInvalidUUID   = errors.New("invalid UUID format")
)

// NewKUID generates a new random KUID
func NewKUID() (*KUID, error) {
	var buf [16]byte
	_, err := rand.Read(buf[:])
	if err != nil {
		return nil, err
	}

	msb := binary.BigEndian.Uint64(buf[0:8])
	lsb := binary.BigEndian.Uint64(buf[8:16])

	return &KUID{msb: msb, lsb: lsb}, nil
}

func FromUUID(uuid string) (*KUID, error) {
	// Validate UUID format with hyphens
	if !strings.HasPrefix(uuid, "") && len(uuid) != 36 {
		return nil, ErrInvalidUUID
	}

	// Check hyphen positions
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return nil, ErrInvalidUUID
	}

	// Remove hyphens and validate hex
	clean := strings.ReplaceAll(uuid, "-", "")
	if len(clean) != 32 {
		return nil, ErrInvalidUUID
	}

	// Validate hex characters
	for _, c := range clean {
		if !strings.ContainsRune("0123456789abcdefABCDEF", c) {
			return nil, ErrInvalidUUID
		}
	}

	// Decode hex string to bytes
	bytes, err := hex.DecodeString(clean)
	if err != nil {
		return nil, ErrInvalidUUID
	}

	return FromBytes(bytes)
}

func FromBytes(b []byte) (*KUID, error) {
	if len(b) != 16 {
		return nil, errors.New("byte slice must be exactly 16 bytes")
	}

	msb := binary.BigEndian.Uint64(b[0:8])
	lsb := binary.BigEndian.Uint64(b[8:16])

	return &KUID{msb: msb, lsb: lsb}, nil
}

// encodeLong encodes a uint64 to base62 in a consistent way
func encodeLong(value uint64) string {
	bytes := make([]byte, size)
	for i := size - 1; i >= 0; i-- {
		bytes[i] = base62Chars[value%base]
		value /= base
	}
	return string(bytes)
}

// decodeLong decodes a base62 string back to uint64
func decodeLong(s string) (uint64, error) {
	if len(s) != size {
		return 0, ErrInvalidLength
	}

	var value uint64
	for i := 0; i < len(s); i++ {
		digit := strings.IndexByte(base62Chars, s[i])
		if digit < 0 {
			return 0, ErrInvalidChar
		}
		value = value*base + uint64(digit)
	}
	return value, nil
}

// String returns the base62 encoded representation of the KUID
func (k KUID) String() string {
	return encodeLong(k.msb) + encodeLong(k.lsb)
}

// FromString creates a KUID from its string representation
func FromString(s string) (*KUID, error) {
	if len(s) != size*2 {
		return nil, ErrInvalidLength
	}

	msb, err := decodeLong(s[:size])
	if err != nil {
		return nil, err
	}

	lsb, err := decodeLong(s[size:])
	if err != nil {
		return nil, err
	}

	return &KUID{msb: msb, lsb: lsb}, nil
}

// Bytes returns the KUID as a 16-byte slice
// ToUUID converts the KUID back to a UUID string format
func (k *KUID) ToUUID() string {
	bytes := k.Bytes()
	uuid := hex.EncodeToString(bytes)

	// Insert hyphens in UUID format: 8-4-4-4-12
	return uuid[0:8] + "-" +
		uuid[8:12] + "-" +
		uuid[12:16] + "-" +
		uuid[16:20] + "-" +
		uuid[20:]
}

func (k *KUID) Bytes() []byte {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[0:8], k.msb)
	binary.BigEndian.PutUint64(b[8:16], k.lsb)
	return b
}

// Equal returns true if two KUIDs are equal
func (k *KUID) Equal(other *KUID) bool {
	if other == nil {
		return false
	}
	return k.msb == other.msb && k.lsb == other.lsb
}
