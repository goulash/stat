// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math"
	"math/rand"
)

// UniformDiscrete gives a discrete uniform distribution between a and b.
type UniformDiscrete struct {
	r *rand.Rand
	a int64
	b int64
}

// NewUniformDiscrete returns a discrete uniform distribution in [a, b).
// Note: this will not work if (b-a) is greater than 2**63.
func NewUniformDiscrete(s rand.Source, a, b int64) *UniformDiscrete {
	if a == b || s == nil {
		panic("invalid input")
	} else if b < a {
		a, b = b, a
	}
	return &UniformDiscrete{rand.New(s), a, b}
}

func (u *UniformDiscrete) String() string {
	return fmt.Sprintf("discrete uniform [%v %v]", u.a, u.b)
}

func (u *UniformDiscrete) Int63() int64 {
	return u.r.Int63n(u.b-u.a) + u.a
}

func (u *UniformDiscrete) P(x int64) (p float64) {
	if x < u.a {
		return 0
	} else if x > u.b {
		return 1
	}
	return float64(x-u.a) / float64(u.b-u.a)
}

func (u *UniformDiscrete) Q(p float64) (x float64) {
	if p < 0 {
		return float64(u.a)
	} else if p > 1 {
		return float64(u.b)
	}
	return math.Floor(p*float64(u.b-u.a) + float64(u.a))
}

func (u *UniformDiscrete) Mean() float64 {
	return float64(u.b-u.a)/2 + float64(u.a)
}
