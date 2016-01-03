// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package test provides statistical tests, such as the ChiSquaredTest.
//
// This package may be removed, as the main function, ChiSquaredTest,
// is only an approximation.
package test

import (
	"fmt"
	"math"
	"os"

	"github.com/goulash/stat"
	"github.com/goulash/stat/dist"
	"github.com/goulash/stat/statutil"
)

// ChiSquaredTest performs a statistical test on the given series.
func ChiSquaredTest(s stat.Series, d dist.Dist, k int, alpha float64) (ok bool, r float64) {
	n := len(s)
	if s == nil || d == nil {
		panic("invalid values submitted")
	}

	// Make a histogram and get values please. We use the distribution inverse for this.
	ks := make([]float64, k+1)
	for i := range ks {
		ks[i] = d.Q(float64(i) / float64(k))
	}
	bs := Bins(ks, s)
	exp := float64(n) / float64(k)

	sum := 0
	var chi2 float64
	for _, b := range bs {
		sum += b
		chi2 += math.Pow(float64(b)-exp, 2) / exp
	}
	if sum != n {
		fmt.Fprintf(os.Stderr, "Error: expected sum to equal %d, got %d.\n", n, sum)
	}
	return chi2 < ChiSquared.Value(float64(k-1), alpha), chi2
}

var ChiSquared = statutil.NewApproxTable(
	[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 40, 50, 60, 70, 80, 90, 100},
	[]float64{0.995, 0.99, 0.975, 0.95, 0.90, 0.10, 0.05, 0.025, 0.01, 0.005},
	[][]float64{
		[]float64{0.000, 0.000, 0.001, 0.004, 0.016, 2.706, 3.841, 5.024, 6.635, 7.879},
		[]float64{0.010, 0.020, 0.051, 0.103, 0.211, 4.605, 5.991, 7.378, 9.210, 10.597},
		[]float64{0.072, 0.115, 0.216, 0.352, 0.584, 6.251, 7.815, 9.348, 11.345, 12.838},
		[]float64{0.207, 0.297, 0.484, 0.711, 1.064, 7.779, 9.488, 11.143, 13.277, 14.860},
		[]float64{0.412, 0.554, 0.831, 1.145, 1.610, 9.236, 11.070, 12.833, 15.086, 16.750},
		[]float64{0.676, 0.872, 1.237, 1.635, 2.204, 10.645, 12.592, 14.449, 16.812, 18.548},
		[]float64{0.989, 1.239, 1.690, 2.167, 2.833, 12.017, 14.067, 16.013, 18.475, 20.278},
		[]float64{1.344, 1.646, 2.180, 2.733, 3.490, 13.362, 15.507, 17.535, 20.090, 21.955},
		[]float64{1.735, 2.088, 2.700, 3.325, 4.168, 14.684, 16.919, 19.023, 21.666, 23.589},
		[]float64{2.156, 2.558, 3.247, 3.940, 4.865, 15.987, 18.307, 20.483, 23.209, 25.188},
		[]float64{2.603, 3.053, 3.816, 4.575, 5.578, 17.275, 19.675, 21.920, 24.725, 26.757},
		[]float64{3.074, 3.571, 4.404, 5.226, 6.304, 18.549, 21.026, 23.337, 26.217, 28.300},
		[]float64{3.565, 4.107, 5.009, 5.892, 7.042, 19.812, 22.362, 24.736, 27.688, 29.819},
		[]float64{4.075, 4.660, 5.629, 6.571, 7.790, 21.064, 23.685, 26.119, 29.141, 31.319},
		[]float64{4.601, 5.229, 6.262, 7.261, 8.547, 22.307, 24.996, 27.488, 30.578, 32.801},
		[]float64{5.142, 5.812, 6.908, 7.962, 9.312, 23.542, 26.296, 28.845, 32.000, 34.267},
		[]float64{5.697, 6.408, 7.564, 8.672, 10.085, 24.769, 27.587, 30.191, 33.409, 35.718},
		[]float64{6.265, 7.015, 8.231, 9.390, 10.865, 25.989, 28.869, 31.526, 34.805, 37.156},
		[]float64{6.844, 7.633, 8.907, 10.117, 11.651, 27.204, 30.144, 32.852, 36.191, 38.582},
		[]float64{7.434, 8.260, 9.591, 10.851, 12.443, 28.412, 31.410, 34.170, 37.566, 39.997},
		[]float64{8.034, 8.897, 10.283, 11.591, 13.240, 29.615, 32.671, 35.479, 38.932, 41.401},
		[]float64{8.643, 9.542, 10.982, 12.338, 14.041, 30.813, 33.924, 36.781, 40.289, 42.796},
		[]float64{9.260, 10.196, 11.689, 13.091, 14.848, 32.007, 35.172, 38.076, 41.638, 44.181},
		[]float64{9.886, 10.856, 12.401, 13.848, 15.659, 33.196, 36.415, 39.364, 42.980, 45.559},
		[]float64{10.520, 11.524, 13.120, 14.611, 16.473, 34.382, 37.652, 40.646, 44.314, 46.928},
		[]float64{11.160, 12.198, 13.844, 15.379, 17.292, 35.563, 38.885, 41.923, 45.642, 48.290},
		[]float64{11.808, 12.879, 14.573, 16.151, 18.114, 36.741, 40.113, 43.195, 46.963, 49.645},
		[]float64{12.461, 13.565, 15.308, 16.928, 18.939, 37.916, 41.337, 44.461, 48.278, 50.993},
		[]float64{13.121, 14.256, 16.047, 17.708, 19.768, 39.087, 42.557, 45.722, 49.588, 52.336},
		[]float64{13.787, 14.953, 16.791, 18.493, 20.599, 40.256, 43.773, 46.979, 50.892, 53.672},
		[]float64{20.707, 22.164, 24.433, 26.509, 29.051, 51.805, 55.758, 59.342, 63.691, 66.766},
		[]float64{27.991, 29.707, 32.357, 34.764, 37.689, 63.167, 67.505, 71.420, 76.154, 79.490},
		[]float64{35.534, 37.485, 40.482, 43.188, 46.459, 74.397, 79.082, 83.298, 88.379, 91.952},
		[]float64{43.275, 45.442, 48.758, 51.739, 55.329, 85.527, 90.531, 95.023, 100.425, 104.215},
		[]float64{51.172, 53.540, 57.153, 60.391, 64.278, 96.578, 101.879, 106.629, 112.329, 116.321},
		[]float64{59.196, 61.754, 65.647, 69.126, 73.291, 107.565, 113.145, 118.136, 124.116, 128.299},
		[]float64{67.328, 70.065, 74.222, 77.929, 82.358, 118.498, 124.342, 129.561, 135.807, 140.169},
	},
)

// Bins returns the number of values of s that land in the len(ks)-1 bins.
func Bins(ks []float64, s stat.Series) []int {
	find := func(x float64) int {
		for i, v := range ks {
			if x < v {
				return i - 1
			}
		}
		return -1
	}

	counts := make([]int, len(ks)-1)
	for _, x := range s {
		i := find(x)
		if i < 0 {
			continue
		}
		counts[i]++
	}
	return counts
}
