// Package bit is home to the bitrange datastructure which is basically an
// ordered list of bits of finite length. It is a packed format.
package bit

import (
	"encoding/json"
	"errors"
	"strconv"
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

// IsSet returns true if the bit located at the given position in the bitrange is set.
// If the position is out-of-bonds for the bitrange, it panics.
// The position starts from 1 (first element).
func (o *Range) IsSet(position int) bool {
	if position > o.size || position < 1 {
		panic("position out-of-bounds: " + strconv.Itoa(position) + " not between 1 and " + strconv.Itoa(o.size))
	}
	// Let's find where the bit is stored
	index := position - 1
	r := index % 64
	l := (index - r) / 64
	if l > len(o.Array)-1 && r != 0 {
		// the bit we are looking for is located in the Leftover
		return (o.Leftover & (1 << r)) != 0
	}
	// the bit is stored in o.Array[l]
	return (o.Array[l] & (1 << r)) != 0
}

// Zeroes returns the unpacked representation of the bitrange as a
// slice of integers. Each integer is the position of an unset bit in the range.
// Finally, the total number of zeroes is returned.
func (o *Range) Zeroes() (list []int, count int) {
	list = make([]int, 0, o.size)
	for k := 0; k < len(o.Array); k++ {
		n := o.Array[k]
		for i := 0; i < 64; i++ {
			if (n & (1 << (i))) == 0 {
				list = append(list, i+1+k*64)
			}
		}
	}

	for i := 0; i < o.LeftoverBitCount; i++ {
		if (o.Leftover & (1 << i)) == 0 {
			list = append(list, i+1+len(o.Array)*64)
		}
	}

	return list, o.size - o.count
}

func (o *Range) Marshal() ([]byte, error) {
	return json.Marshal(*o)
}

func (o *Range) UnMarshal(data []byte) error {
	return json.Unmarshal(data, o)
}
