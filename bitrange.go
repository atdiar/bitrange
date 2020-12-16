// Package bit is home to the bitrange datastructure which is basically an
// ordered list of bit of finite length. It is a packed format.
package bit

import (
	"encoding/json"
	"errors"
)

// Range defines the structure of a bitrange
type Range struct {
	Array            []uint64
	Leftover         uint64
	LeftoverBitCount int
	count            int
	size             int
}

// NewRange returns a new bitrange datastructure
func NewRange(size int) *Range {
	r := (size) % 64
	l := (size - r) / 64
	if r == 0 {
		l--
	}

	overflow := size - l*64

	return &Range{make([]uint64, l), 0, overflow, 0, size}
}

// Set will set the bit at the given position in the bitrange,
// starting from 1. (the first bit)
func (o *Range) Set(position int) error {
	if position > o.size || position == 0 {
		return errors.New("cannot set bit at given position. Out-of-range.")
	}
	r := position % 64
	q := (position - r) / 64
	if r == 0 {
		q--
		if q < 0 {
			q = 0
		}
	}

	overflow := position - q*64

	pos := position - q*64

	if q == len(o.Array) && overflow != 0 {
		o.Leftover |= o.Leftover ^ (1 << (pos - 1))
		o.count++
		return nil
	}

	o.Array[q] |= o.Array[q] ^ (1 << (pos - 1))
	o.count++
	return nil
}

// Zeroes returns the number of bits that have been unset as well as their
// position in the bitrange.
func (o *Range) Zeroes() (list []int, count int) {
	list = make([]int, 0, o.size)
	for k := 0; k < len(o.Array); k++ {
		n := o.Array[k]
		for i := 0; i < 64; i++ {
			if (n & (1 << (i))) == 0 {
				list = append(list, i+1+k*64)
				count++
			}
		}
	}
	if o.Leftover != 0 {
		for i := 0; i < o.LeftoverBitCount; i++ {
			if (o.Leftover & (1 << i)) == 0 {
				list = append(list, i+1+len(o.Array)*64)
				count++
			}
		}
	}
	return list, count
}

func (o *Range) Marshal() ([]byte, error) {
	return json.Marshal(*o)
}

func (o *Range) UnMarshall(data []byte) error {
	return json.Unmarshal(data, o)
}
