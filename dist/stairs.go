// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math/rand"
)

// Stairs returns the index of the first probability value that exceeds the
// random value between 0.0 and 1.0.
//
// For example, given the following list:
//
//  [0.0, 0.3, 0.6, 0.6, 0.9]
//
// We can imagine the following stairs, with indices from 0 to 4.
//
//                      .
//            .____.____|
//       .____|
//  .____|
//  0   0.3  0.6  0.6  0.9
//  0    1    2    3    4
//
// The probability for index 0 and 3 is 0.0. The probability for the rest is is
// the value minus the previous divided by the last value in the list, in this
// case 0.9. The indices 1, 2, and 4 all have the same probability then.
// Therefore, we get the same result if we pass in the list:
//
//  [0, 3, 6, 6, 9]
//
// The only requirement on the numbers in the list are that they are monotonically
// increasing. Failing this requirement will cause NewStairs to panic.
type Stairs struct {
	r *rand.Rand
	p []float64
	z int64
}

func NewStairs(s rand.Source, p ...float64) *Stairs {
	n := len(p)
	if n == 0 {
		panic("list of probabilities cannot be zero")
	}

	prev := 0.0
	denom := p[n-1]
	xps := make([]float64, n)
	for i, r := range p {
		if r < prev {
			panic("probabilities must be monotonically increasing")
		}
		prev = r
		xps[i] = r / denom
	}

	return &Stairs{
		r: rand.New(s),
		p: xps,
		z: int64(n - 1),
	}
}

func (s *Stairs) String() string {
	return fmt.Sprintf("discrete stairs %v", s.p)
}

func (s *Stairs) Int63() int64 {
	p := s.r.Float64()
	for i, r := range s.p {
		if r > p {
			return int64(i)
		}
	}
	return s.z
}

func (s *Stairs) P(x int64) (p float64) {
	if x < 0 || x > s.z {
		return 0.0
	}
	return s.p[x] - s.p[x-1]
}

func (s *Stairs) Q(p float64) (x float64) {
	if p < 0.0 {
		return 0.0
	} else if p >= 1.0 {
		return float64(s.z)
	}
	for i, r := range s.p {
		if r > p {
			return float64(i)
		}
	}
	return float64(s.z)
}

func (s *Stairs) Mean() float64 {
	var mean float64
	for i, p := range s.p[1:] {
		mean += float64(i+1) * (p - s.p[i])
	}
	return mean
}
