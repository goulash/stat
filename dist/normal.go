// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math/rand"
)

// Normal distribution with mean and standard deviation.
type Normal struct {
	r    *rand.Rand
	mean float64
	std  float64
}

func NewNormal(s rand.Source, mean, std float64) *Normal {
	if s == nil {
		panic("s cannot be nil")
	}

	return &Normal{rand.New(s), mean, std}
}

func (n *Normal) String() string {
	return fmt.Sprintf("normal [μ=%f σ=%f]", n.mean, n.std)
}

func (n *Normal) Float64() float64 {
	return n.r.NormFloat64()*n.std + n.mean
}

func (n *Normal) Mean() float64 { return n.mean }
func (n *Normal) Var() float64  { return n.std * n.std }
func (n *Normal) Std() float64  { return n.std }

func (n *Normal) Z(x float64) float64 { return (x - n.mean) / n.std }
