// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"math/rand"
	"testing"
)

func TestStairsMean(z *testing.T) {
	type test struct {
		P []float64
		M float64
	}
	tests := []test{
		test{[]float64{0.0, 0.3, 0.6, 0.6, 0.9}, (0.3/0.9*1 + 0.3/0.9*2 + 0.3/0.9*4)},
	}

	source := rand.NewSource(0)
	for _, t := range tests {
		s := NewStairs(source, t.P...)
		if s.Mean() != t.M {
			z.Errorf("Stairs%v.Mean() = %f, want %f", t.P, s.Mean(), t.M)
		}
	}
}
