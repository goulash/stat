// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package stat

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

var NaN = math.NaN()

type test struct {
	A, B             Series
	Min, Max         float64
	Mean, Median     float64
	Var, VarP        float64
	Std, StdP        float64
	Autocov, Autocor []float64 // start with lag 0
	Cov, CovP        float64
	Cor              float64
}

func TestString(z *testing.T) {
	assert := assert.New(z)

	t := Series{1, 2, 3, 4, 5}
	assert.Equal("[1 2 3 4 5]", t.String(), "string output should be the same")

	t.Reset()
	assert.Equal(0, t.Len(), "after reset, series should be empty")
	assert.Equal("[]", t.String(), "this is what an empty series should look like")
}

func TestHeadTail(z *testing.T) {
	assert := assert.New(z)
	t := Series{1, 2, 3, 4, 5}
	assert.Equal("[1 2 3]", t.Head(3).String(), "first three elements")
	assert.Equal("[3 4 5]", t.Tail(3).String(), "last three elements")
	assert.Equal("[1]", t.Head(1).String(), "first element")
	assert.Equal("[5]", t.Tail(1).String(), "last element")
	assert.Equal("[]", t.Head(0).String(), "want nothing, get nothing")
	assert.Equal(Series{}, t.Tail(0), "want nothing, get nothing")
	assert.Equal(t.Tail(0), t.Head(0), "technically the same")
	assert.Equal(t.Tail(5), t.Head(5), "technically the same")
	assert.Panics(func() { t.Head(-1) }, "expect out-of-bounds panic")
	assert.Panics(func() { t.Head(6) }, "expect out-of-bounds panic")
}

func TestMedian(z *testing.T) {
	assert := assert.New(z)
	t := Series{1, 2, 3, 4, 5}
	assert.Equal(3.0, t.Median(), "median should be equal")
	t.Append(6)
	assert.Equal((3.0+4.0)/2.0, t.Median(), "median should be equal")
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func fix(f float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(f*output)) / output
}

func rAssert(z *testing.T) func(float64, float64, string) {
	p := 7 // precision
	return func(want, got float64, msg string) {
		if math.IsNaN(want) {
			assert.NotEqual(z, got, got, "should be NaN: %s", msg)
		} else {
			assert.Equal(z, fix(want, p), fix(got, p), msg)
		}
	}
}

func TestSeries(z *testing.T) {
	assert := rAssert(z)
	for _, t := range tests {
		a, b := t.A, t.B
		assert(t.Min, a.Min(), "min should be equal")
		assert(t.Max, a.Max(), "max should be equal")
		assert(t.Mean, a.Mean(), "mean should be equal")
		assert(t.Median, a.Median(), "median should be equal")
		assert(t.Var, a.Var(), "var should be equal")
		assert(t.VarP, a.VarP(), "varp should be equal")
		assert(t.Std, a.Std(), "std should be equal")
		assert(t.StdP, a.StdP(), "stdp shoulb be equal")
		for i, x := range t.Autocov {
			assert(x, a.Autocov(i), fmt.Sprintf("autocov with lag %d should be equal", i))
		}
		for i, x := range t.Autocor {
			assert(x, a.Autocor(i), fmt.Sprintf("autocorr with lag %d should be equal", i))
		}

		assert(t.Cov, a.Cov(b), "cov should be equal")
		assert(t.CovP, a.CovP(b), "covp should be equal")
		assert(t.Cor, a.Cor(b), "cor should be equal")
	}
}

var tests = []test{
	{
		A: Series{
			0.38809179, 0.94113008, 0.15350705, 0.03311646, 0.68168087, 0.21719990,
			0.32123922, 0.57085251, 0.53576882, 0.38965630, 0.27487263, 0.90783122,
		},
		B: Series{
			0.53644778, 0.51250793, 0.47027137, 0.39856790, 0.87766664, 0.48159821,
			0.83109965, 0.07217068, 0.62203444, 0.21900594, 0.30220383, 0.83902456,
		},
		Min:    0.03311646,
		Max:    0.94113008,
		Mean:   0.45124557083333333,
		Median: 0.388874045000000,
		Var:    0.0815505430886907,
		VarP:   0.0747546644979665,
		Std:    0.285570557111007,
		StdP:   0.273412992555157,
		Autocov: []float64{
			0.081550543, -0.024985019, -0.032008039, 0.036685824, -0.002588976, -0.024764554,
			-0.013166106, 0.059028468, -0.050665201, -0.113928076, 0.175025168, NaN, NaN,
		},
		Autocor: []float64{
			1.0, -0.32315216, -0.44856680, 0.48593907, -0.03693304, -0.32909314,
			-0.16419487, 0.65930820, -0.45680962, -0.83541672, 1.0, NaN, NaN,
		},
		Cov:  0.0245006797759880,
		CovP: 0.0224589564613224,
		Cor:  0.341571101278264,
	},
}
