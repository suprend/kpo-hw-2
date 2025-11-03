package id

import (
	"crypto/rand"
	"math/big"
	"sync"
	"time"

	"kpo-hw-2/internal/domain"
)

type ULIDGenerator struct {
	mu sync.Mutex
}

func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

func (g *ULIDGenerator) NewID() (domain.ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	ts := uint64(time.Now().UnixMilli())
	var randomness [10]byte
	if _, err := rand.Read(randomness[:]); err != nil {
		return "", err
	}

	raw := encodeULID(ts, randomness)
	return domain.ParseID(raw)
}

var _ domain.IDGenerator = (*ULIDGenerator)(nil)

func encodeULID(ts uint64, randomness [10]byte) string {
	var data [16]byte
	data[0] = byte(ts >> 40)
	data[1] = byte(ts >> 32)
	data[2] = byte(ts >> 24)
	data[3] = byte(ts >> 16)
	data[4] = byte(ts >> 8)
	data[5] = byte(ts)
	copy(data[6:], randomness[:])

	return encodeBase32(data[:])
}

func encodeBase32(data []byte) string {
	n := new(big.Int).SetBytes(data)
	base := big.NewInt(32)
	zero := big.NewInt(0)

	var chars [26]byte
	for i := len(chars) - 1; i >= 0; i-- {
		if n.Cmp(zero) <= 0 {
			chars[i] = domain.ULIDAlphabet[0]
			continue
		}

		mod := new(big.Int)
		n.DivMod(n, base, mod)
		chars[i] = domain.ULIDAlphabet[mod.Int64()]
	}

	return string(chars[:])
}
