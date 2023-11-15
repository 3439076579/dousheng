package bloom

import (
	"context"
	"github.com/spaolacci/murmur3"
	"math"
)
import "github.com/go-redis/redis/v8"

/*
	BloomFilter

Calculate bitLen from probably of a false positive and ExceptData
BitLen = ExceptData*In(Rate)/In2*In2

Calculate hashFunc from BitLen and ExceptData
hashFunc=BitLen/ExceptData*In(2)
*/
const (
	key          = "bloom_filter"
	maxBitMapLen = 512 * 1024 * 1024 * 8
)

func calBitLen(ExceptData uint64, rate float64) uint64 {
	var bitLen uint64
	bitLen = -(uint64)((float64)(ExceptData) * math.Log(rate) * math.Log2E * math.Log2E)
	// if bitLen > MaxBitMapLen
	if bitLen > maxBitMapLen {
		bitLen = maxBitMapLen
	}
	return bitLen
}

func calHashFunCount(ExceptData uint64, BitLen uint64) uint64 {
	count := uint64((float64)(BitLen) / (float64)(ExceptData) * math.Ln2)
	if count == 0 {
		panic("ExceptData is too large")
	}
	return count
}

type BloomFilter struct {
	// key is name of BloomFilter,default value is "bloom_filter"
	key string
	// bitLen is length of BloomFilter
	BitLen    uint64
	MapCount  uint64
	DataCount uint64
	MJRate    float64
	store     *redis.Client
}

func NewBloomFilter(ExceptData uint64, Rate float64, store *redis.Client, name ...string) *BloomFilter {
	var Key string = key
	if len(name) != 0 {
		Key = name[0]
	}
	bf := &BloomFilter{
		key:       Key,
		BitLen:    calBitLen(ExceptData, Rate),
		DataCount: 0,
		store:     store,
	}
	bf.MapCount = calHashFunCount(ExceptData, bf.BitLen)
	return bf
}

func (bf *BloomFilter) IsExist(key string) (bool, error) {
	offsets := bf.getOffset(key)
	return bf.getBit(offsets)
}
func (bf *BloomFilter) getOffset(key string) []uint64 {
	var offsets []uint64
	var data = ([]byte)(key)
	for i := uint64(0); i < bf.MapCount; i++ {
		data = append(data, byte(i))
		sum := murmur3.Sum64(data)
		offsets = append(offsets, sum%bf.BitLen)
	}
	return offsets
}

func (bf *BloomFilter) Add(ctx context.Context, key string) error {
	offsets := bf.getOffset(key)
	err := bf.setBit(ctx, offsets)
	if err != nil {
		return err
	}
	return nil
}
