// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package statutil

import "math"

// ApproxTable is a distribution table approximation.
//
// In particular, it is good for complex distributions that I do not want to implement,
// such as the chi squared distribution.
type ApproxTable struct {
	rows  []float64
	cols  []float64
	table [][]float64
}

func NewApproxTable(rows, cols []float64, table [][]float64) *ApproxTable {
	// Make sure that the input is valid.
	n, m := len(rows), len(cols)
	if n == 0 || m == 0 {
		panic("require valid parameters")
	}
	if len(table) != n {
		panic("table does not have n rows")
	}
	for _, xs := range table {
		if len(xs) != m {
			panic("table row does not have m columns")
		}
	}

	return &ApproxTable{
		rows:  rows,
		cols:  cols,
		table: table,
	}
}

func (at ApproxTable) Value(r, c float64) float64 {
	findClosest := func(xs []float64, y float64) int {
		md := math.Abs(xs[0] - y)
		mi := 0
		for i, x := range xs {
			d := math.Abs(x - y)
			if d < md {
				md = d
				mi = i
			}
		}
		return mi
	}

	ri := findClosest(at.rows, r)
	ci := findClosest(at.cols, c)
	return at.table[ri][ci]
}
