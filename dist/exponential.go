// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math"
	"math/rand"
)

// Exponential distribution with rate of arrival.
type Exponential struct {
	r      *rand.Rand
	lambda float64
}

func NewExponential(lambda float64, s rand.Source) *Exponential {
	if lambda <= 0 || s == nil {
		panic("lambda must be positive")
	}

	return &Exponential{rand.New(s), lambda}
}

func (e Exponential) String() string {
	return fmt.Sprintf("exponential [%v]", e.lambda)
}

func (e Exponential) Float64() float64 {
	// Inversion method works here.
	return e.Q(e.r.Float64())
}

func (e Exponential) P(x float64) float64 {
	if x < 0 {
		return 0
	}

	return 1 - math.Exp(-e.lambda*x)
}

func (e Exponential) Q(p float64) float64 {
	if p < 0 {
		return 0
	} else if p == 1 {
		return Inf
	}

	return -math.Log(1-p) / e.lambda
}

func (e Exponential) Mean() float64 {
	return 1 / e.lambda
}
