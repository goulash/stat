// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math/rand"
)

// Uniform gives a uniform distribution between a and b.
type Uniform struct {
	r *rand.Rand
	a float64
	b float64
}

func NewUniform(s rand.Source, a, b float64) *Uniform {
	if a == b || s == nil {
		panic("invalid input")
	} else if b < a {
		a, b = b, a
	}
	return &Uniform{rand.New(s), a, b}
}

func (u *Uniform) String() string {
	return fmt.Sprintf("uniform [%v %v]", u.a, u.b)
}

func (u *Uniform) Float64() float64 {
	// Inversion method works here.
	return u.Q(u.r.Float64())
}

func (u *Uniform) P(x float64) (p float64) {
	if x < u.a {
		return 0
	} else if x > u.b {
		return 1
	}
	return (x - u.a) / (u.b - u.a)
}

func (u *Uniform) Q(p float64) (x float64) {
	if p < 0 {
		return u.a
	} else if p > 1 {
		return u.b
	}
	return p*(u.b-u.a) + u.a
}

func (u *Uniform) Mean() float64 {
	return u.Q(0.5)
}
