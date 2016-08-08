// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

type rounded struct {
	c Continuous
	f func(float64) int64
}

func (d rounded) Int63() int64 {
	return d.f(d.c.Float64())
}

func Floor(c Continuous) Discrete {
	return rounded{c, func(f float64) int64 {
		return int64(f)
	}}
}

func Ceil(c Continuous) Discrete {
	return rounded{c, func(f float64) int64 {
		n := int64(f)
		if float64(n) == f {
			return n
		} else {
			return n + 1
		}
	}}
}

func Round(c Continuous) Discrete {
	return rounded{c, func(f float64) int64 {
		return int64(f + 0.5)
	}}
}
