// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package dist provides statistical probability distributions for random variables.
package dist

import "github.com/goulash/stat"

// Nabbed these from math...
const (
	NaN    = 0x7FF8000000000001
	Inf    = 0x7FF0000000000000
	NegInf = 0xFFF0000000000000
)

type DistP interface {
	// Name of the distribution
	String() string

	// Float64 returns a random value from the distribution.
	Float64() float64

	// P returns the probability that a value from the distribution is less than x.
	// That is, P represents the cumulative probability function of the distribution.
	P(x float64) (p float64)
}

// Dist is implemented by all distributions that give the inverse CDF function.
// Because this is not possible for all distributions, it remains optional.
type Dist interface {
	DistP

	// Q returns the p-quantile of the distribution, this is the inverse CDF function.
	Q(p float64) (x float64)
}

// PDF returns the probability that a value lands between a and b, where a <= b.
func PDF(d DistP, a, b float64) (p float64) {
	if b < a {
		return -1
	}
	return d.P(b) - d.P(a)
}

func GenSeries(d DistP, n int) stat.Series {
	if n <= 0 {
		return nil
	}
	s := make(stat.Series, n)
	for i := range s {
		s[i] = d.Float64()
	}
	return s
}
