// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

type pass struct {
	f func() float64
}

func (p pass) Float64() float64 { return p.f() }

func Lowpass(c Continuous, high float64) Continuous {
	return &pass{func() float64 {
		x := c.Float64()
		for high < x {
			x = c.Float64()
		}
		return x
	}}
}

func Highpass(c Continuous, low float64) Continuous {
	return &pass{func() float64 {
		x := c.Float64()
		for x < low {
			x = c.Float64()
		}
		return x
	}}
}

func Midpass(c Continuous, low, high float64) Continuous {
	return &pass{func() float64 {
		x := c.Float64()
		for x < low || high < x {
			x = c.Float64()
		}
		return x
	}}
}
