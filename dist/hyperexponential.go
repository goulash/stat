// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math/rand"
)

// HyperExponential distribution with k rates of arrival.
type HyperExponential struct {
	r       *rand.Rand
	stairs  *Stairs
	lambdas []float64
}

func NewHyperExponential(s rand.Source, probs, lambdas []float64) *HyperExponential {
	if s == nil {
		panic("random source cannot be zero")
	}
	if len(probs) != len(lambdas) {
		panic("list of probabilities must have same length as lambdas")
	}
	for _, l := range lambdas {
		if l <= 0 {
			panic("lambda must be positive")
		}
	}

	stairs := NewStairs(s, probs...)
	return &HyperExponential{
		r:       stairs.r,
		stairs:  stairs,
		lambdas: lambdas,
	}
}

func (e *HyperExponential) String() string {
	return fmt.Sprintf("hyper-exponential %v", e.lambdas)
}

func (e *HyperExponential) Float64() float64 {
	return e.r.ExpFloat64() / e.lambdas[e.stairs.Int63()]
}

func (e *HyperExponential) Mean() float64 {
	var mean, prev float64
	for i, p := range e.stairs.p {
		mean += (p - prev) / e.lambdas[i]
		prev = p
	}
	return mean
}

func (e *HyperExponential) SecondMoment() float64 {
	var mean, prev float64
	for i, p := range e.stairs.p {
		mean += (2 * (p - prev)) / (e.lambdas[i] * e.lambdas[i])
		prev = p
	}
	return mean

}

func (e *HyperExponential) Var() float64 {
	return e.SecondMoment() - (e.Mean() * e.Mean())
}
