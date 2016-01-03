// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// +build ignore

package stat

import "math"

// RunAll calculates the running means, variances, and standard deviation over time.
type RunAll struct {
	z Run

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
