// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math"
	"math/rand"
)

// Poisson returns random numbers according to a poisson distribution.
//
// The poisson distribution models the number of arrivals in a set time interval
// when the inter-arrival times are exponentially distributed.
// This is why
type Poisson struct {
	r      *rand.Rand
	lambda float64
}

func NewPoisson(s rand.Source, lambda float64) *Poisson {
	if lambda <= 0 || s == nil {
		panic("lambda must be positive")
	}

	return &Poisson{rand.New(s), lambda}
}

func (p Poisson) String() string {
	return fmt.Sprintf("poisson [%v]", p.lambda)
}

func (p Poisson) Int63() int64 {
	a := math.Exp(-p.lambda)
	var b float64 = 1
	var k int64 = -1
	for b <= a {
		b *= p.r.Float64()
		k++
	}
	return k
}

func (p Poisson) Mean() float64 {
	return p.lambda
}

func (p Poisson) Var() float64 {
	return p.lambda
}
