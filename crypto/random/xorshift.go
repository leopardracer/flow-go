package random

import (
	"encoding/binary"
	"fmt"
)

// xorshifts is a set of xorshift128+ pseudo random number generators
// each xorshift+ has a 128-bits state
// using a set of xorshift128+ allows initializing the set with a larger
// seed than 128 bits
type xorshifts struct {
	states     []xorshiftp
	stateIndex int
}

// xorshiftp is a single xorshift128+ PRG
// the internal state is 128 bits
// http://vigna.di.unimi.it/ftp/papers/xorshiftplus.pdf
type xorshiftp struct {
	a, b uint64
}

// NewRand returns a new set of xorshift128+ PRGs
// the length of the seed fixes the number of xorshift128+ to initialize:
// each 16 bytes of the seed initilize an xorshift128+ instance
// To make sure the seed entropy is optimal, the function checks that len(seed)
// is a multiple 16 (PRG state size)
func NewRand(seed []byte) (*xorshifts, error) {
	// safety check
	if len(seed) == 0 || len(seed)%16 != 0 {
		return nil, fmt.Errorf("seed length should be a non-zero multiple of 16")
	}
	// create the xorshift128+ instances
	states := make([]xorshiftp, 0, len(seed)/16)
	// initialize the xorshift128+ with the seed
	for i := 0; i < cap(states); i++ {
		states = append(states, xorshiftp{
			a: binary.BigEndian.Uint64(seed[i*16 : i*16+8]),
			b: binary.BigEndian.Uint64(seed[i*16+8 : (i+1)*16]),
		})
	}
	// check states are not zeros
	for _, x := range states {
		if x.a|x.b == 0 {
			return nil, fmt.Errorf("the seed of xorshift+ cannot be zero")
		}
	}
	// initial next
	for _, x := range states {
		x.next()
	}
	// init the xorshifts
	rand := &xorshifts{
		states:     states,
		stateIndex: 0,
	}
	return rand, nil
}

// next generates updates the state of a single xorshift128+
func (x *xorshiftp) next() {
	// the xorshift+ shift parameters chosen for this instance
	shifts := []byte{23, 17, 26}
	var tmp uint64 = x.a
	x.a = x.b
	tmp ^= tmp << shifts[0]
	tmp ^= tmp >> shifts[1]
	tmp ^= x.b ^ (x.b >> shifts[2])
	x.b = tmp
}

// prn generated a Pseudo-random number out of the current state of
// xorshift128+ at index (index)
// prn does not change any prg state
func (x *xorshiftp) prn() uint64 {
	return x.a + x.b
}

// IntN returns an uint64 pseudo-random number in [0,n-1]
// using the xorshift+ of the current index. The index is updated
// to use another xorshift+ at the next round
func (x *xorshifts) IntN(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("input must be positive")
	}
	res := x.states[x.stateIndex].prn() % uint64(n)
	// update the state
	x.states[x.stateIndex].next()
	// update the index
	x.stateIndex = (x.stateIndex + 1) % len(x.states)
	return int(res), nil
}

// Permutation returns a permutation of the set [0,n-1]
// it implements Fisher-Yates Shuffle (inside-out variant) using (x) as a random source
// the output space grows very fast with (!n) so that input (n) and the seed length
// (which fixes the internal state length of xorshifts ) should be chosen carefully
// O(n) space and O(n) time
func (x *xorshifts) Permutation(n int) ([]int, error) {
	if n <= 0 {
		return nil, fmt.Errorf("arguments to PermutateSubset must be positive")
	}
	items := make([]int, n)
	for i := 0; i < n; i++ {
		j, _ := x.IntN(i + 1)
		if j != i {
			items[i] = items[j]
		}
		items[j] = i
	}
	return items, nil
}

// SubPermutation returns the m first elements of a permutation of [0,n-1]
// It implements Fisher-Yates Shuffle using x as a source of randoms
// O(n) space and O(n) time
func (x *xorshifts) SubPermutation(n int, m int) ([]int, error) {
	if m <= 0 {
		return nil, fmt.Errorf("arguments to PermutateSubset must be positive")
	}
	items, _ := x.Permutation(n)
	return items[:m], nil
}
