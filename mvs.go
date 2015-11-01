// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package stat

import "math"

// MVS calculates the running mean, variance, and standard deviation.
type MVS struct {
	n int64
	m float64
	s float64

	max float64
	min float64
}

// Reset the values to zero.
func (r *MVS) Reset() {
	r.n = 0
	r.m = 0
	r.s = 0
	r.max = math.Inf(-1)
	r.min = math.Inf(1)
}

// Add a new value to the run.
func (r *MVS) Add(x float64) {
	r.n++
	if r.n == 1 {
		r.m = x
		r.max = x
		r.min = x
		return
	}
	if r.max < x {
		r.max = x
	}
	if r.min > x {
		r.min = x
	}

	m := r.m + (x-r.m)/float64(r.n)
	r.s = r.s + (x-r.m)*(x-m)
	r.m = m
}

func (r *MVS) AddN(n int64, x float64) {
	if r.n == 0 {
		r.n = n
		r.m = x
		r.max = x
		r.min = x
		return
	}
	if r.max < x {
		r.max = x
	}
	if r.min > x {
		r.min = x
	}

	// TODO: Double check this calculation!!!
	i := float64(n - r.n)
	m := r.m + (i*x-i*r.m)/float64(n)
	r.s = r.s + i*(x-r.m)*(x-m)
	r.m = m
	r.n = n
}

// Max returns the max.
func (r MVS) Max() float64 {
	return r.max
}

// Min returns the min.
func (r MVS) Min() float64 {
	return r.min
}

// Mean returns the mean.
func (r MVS) Mean() float64 {
	return r.m
}

// PVar returns the population variance.
func (r MVS) PVar() float64 {
	if r.n <= 1 {
		return 0
	}
	return r.s / float64(r.n)
}

// SVar returns the sample variance.
func (r MVS) SVar() float64 {
	if r.n <= 1 {
		return 0
	}
	return r.s / float64(r.n-1)
}

// PStd returns the population standard deviation.
func (r MVS) PStd() float64 {
	return math.Sqrt(r.PVar())
}

// SStd returns the sample standard deviation.
func (r MVS) SStd() float64 {
	return math.Sqrt(r.SVar())
}

// RunningMVS calculates the running means, variances, and standard deviation over time.
type RunningMVS struct {
	z MVS

	ts []Time
	ms Series
	ss Series
}

// I don't want to take along the baggage of a time.Time.
type Time uint64

func (r *RunningMVS) Add(t Time, x float64) {
	r.z.Add(x)
	r.ts = append(r.ts, t)
	r.ms.Append(r.z.m)
	r.ss.Append(r.z.s)
}

func (r RunningMVS) Times() []Time {
	ts := make([]Time, len(r.ts))
	copy(ts, r.ts)
	return ts
}

func (r RunningMVS) Means() Series {
	ms := make(Series, len(r.ms))
	copy(ms, r.ms)
	return ms
}

func (r RunningMVS) PVars() Series {
	ss := make(Series, len(r.ss))
	n := float64(r.z.n)
	for i, m := range r.ss {
		ss[i] = m / n
	}
	return ss
}

func (r RunningMVS) SVars() Series {
	ss := make(Series, len(r.ss))
	n := float64(r.z.n - 1)
	for i, m := range r.ss {
		ss[i] = m / n
	}
	return ss
}

func (r RunningMVS) PStds() Series {
	ss := make(Series, len(r.ss))
	n := float64(r.z.n)
	for i, m := range r.ss {
		ss[i] = math.Sqrt(m / n)
	}
	return ss
}

func (r RunningMVS) SStds() Series {
	ss := make(Series, len(r.ss))
	n := float64(r.z.n - 1)
	for i, m := range r.ss {
		ss[i] = math.Sqrt(m / n)
	}
	return ss
}
