// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math"
	"math/rand"
)

// LogNormal distribution with mean and standard deviation.
//
// See:
//  http://stackoverflow.com/questions/23699738
//  http://blogs.sas.com/content/iml/2014/06/04/simulate-lognormal-data-with-specified-mean-and-variance.html
type LogNormal struct {
	r    *rand.Rand
	mean float64
	std  float64
}

func NewLogNormal(rs rand.Source, m, s float64) *LogNormal {
	if rs == nil {
		panic("rs cannot be nil")
	}

	// Scale mean and std down to expected mean and std
	m2, s2 := m*m, s*s
	mu := math.Log((m * m) / math.Sqrt(m2+s2))
	sigma := math.Sqrt(math.Log((m2 + s2) / m2))

	return &LogNormal{rand.New(rs), mu, sigma}
}

func (n *LogNormal) String() string {
	return fmt.Sprintf("lognormal [μ=%f σ=%f]", n.mean, n.std)
}

func (n *LogNormal) Float64() float64 {
	return math.Exp(n.r.NormFloat64()*n.std + n.mean)
}

func (n *LogNormal) Mean() float64 { return n.mean }
func (n *LogNormal) Var() float64  { return n.std * n.std }
func (n *LogNormal) Std() float64  { return n.std }

func (n *LogNormal) Z(x float64) float64 { return (x - n.mean) / n.std }
