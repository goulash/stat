// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package stat provides statistic functions and types.
package stat

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
)

// The Series type is a slice of float64 values.
type Series []float64

func (s *Series) Reset()              { *s = make(Series, 0) }
func (s *Series) Append(f ...float64) { *s = append(*s, f...) }
func (s *Series) Append1(f float64)   { *s = append(*s, f) }

func (s Series) Copy() Series {
	t := make(Series, len(s))
	copy(s, t)
	return t
}

func (s Series) WriteFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewWriter(f)
	defer buf.Flush()
	for _, f := range s {
		// TODO: it's probably better to do this with strconv
		_, err := buf.WriteString(fmt.Sprintf("%f\n", f))
		if err != nil {
			return err
		}
	}
	return err
}

func (s Series) String() string {
	if len(s) == 0 {
		return "[]"
	}

	var buf bytes.Buffer
	buf.WriteRune('[')
	buf.WriteString(strconv.FormatFloat(s[0], 'f', -1, 64))
	for _, x := range s[1:] {
		buf.WriteRune(' ')
		buf.WriteString(strconv.FormatFloat(x, 'f', -1, 64))
	}
	buf.WriteRune(']')
	return buf.String()
}

// The following methods are provided for ease-of-use.
// The documentation can be found with the functions of same name.

func (s Series) Len() int                           { return len(s) }
func (s Series) Head(n int) Series                  { return Head(s, n) }
func (s Series) Tail(n int) Series                  { return Tail(s, n) }
func (s Series) Max() float64                       { return Max(s) }
func (s Series) Min() float64                       { return Min(s) }
func (s Series) Mean() float64                      { return Mean(s) }
func (s Series) Median() float64                    { return Median(s) }
func (s Series) SVar() float64                      { return SVar(s) }
func (s Series) PVar() float64                      { return PVar(s) }
func (s Series) SStd() float64                      { return SStd(s) }
func (s Series) PStd() float64                      { return PStd(s) }
func (s Series) SSkew() float64                     { return SSkew(s) }
func (s Series) PSkew() float64                     { return PSkew(s) }
func (s Series) Autocov(lag int) float64            { return Autocov(s, lag) }
func (s Series) Autocorr(lag int) float64           { return Autocorr(s, lag) }
func (s Series) SCovar(t Series) float64            { return SCovar(s, t) }
func (s Series) PCovar(t Series) float64            { return PCovar(s, t) }
func (s Series) SCorr(t Series) float64             { return SCorr(s, t) }
func (s Series) PCorr(t Series) float64             { return PCorr(s, t) }
func (s Series) Map(f func(float64) float64) Series { return Map(s, f) }
func (s Series) Add1(f float64) Series              { return Add1(s, f) }
func (s Series) Mul1(f float64) Series              { return Mul1(s, f) }
func (s Series) Sub1(f float64) Series              { return Sub1(s, f) }
func (s Series) Div1(f float64) Series              { return Div1(s, f) }
func (s Series) Add(t Series) Series                { return Add(s, t) }
func (s Series) Mul(t Series) Series                { return Mul(s, t) }
func (s Series) Sub(t Series) Series                { return Sub(s, t) }
func (s Series) Div(t Series) Series                { return Div(s, t) }

// Head returns the first n values from s.
//
// If n > len(s), an out-of-bounds panic will occur.
// The returned slice is a slice from s, not a new series.
func Head(s Series, n int) Series {
	return s[:n]
}

// Tail returns the last n values from s.
//
// If n > len(s), an out-of-bounds panic will occur.
// The returned slice is a slice from s, not a new series.
func Tail(s Series, n int) Series {
	return s[len(s)-n:]
}

// Max returns the maximum value in the series, or -∞ if the series is empty.
func Max(s Series) float64 {
	if len(s) == 0 {
		return math.Inf(-1)
	}

	m := s[0]
	for _, x := range s[1:] {
		if m < x {
			m = x
		}
	}
	return m
}

// Min returns the minimum value in the series, or +∞ if the series is empty.
func Min(s Series) float64 {
	if len(s) == 0 {
		return math.Inf(1)
	}

	m := s[0]
	for _, x := range s[1:] {
		if m > x {
			m = x
		}
	}
	return m
}

// Mean returns the empirical mean of the series s.
//
// The mean calculated here is the running mean, which ensures that
// an answer is given regardless of how long s is. The accuracy of
// the answer suffers however.
//
// If s is empty, NaN is returned.
func Mean(s Series) float64 {
	if len(s) == 0 {
		return math.NaN()
	}

	var m float64
	for i, x := range s {
		m += (x - m) / float64(i+1)
	}
	return m
}

// Median returns the median of the series s.
//
// If s has an even number of elements, the mean of the two middle
// elements is returned.
//
// If s is empty or has only one element, NaN is returned.
//
// Calculating the median requires sorting a copy of the series.
func Median(s Series) float64 {
	if len(s) <= 1 {
		return math.NaN()
	}

	n := len(s)
	t := s.Copy()
	sort.Float64s(t)
	if n%2 == 0 {
		return (t[n/2] + t[n/2+1]) / 2
	}
	return t[n/2]
}

// SVar returns the sample variance of the series.
//
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func SVar(s Series) float64 {
	return variance(s, true)
}

// PVar returns the population variance of the series.
//
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func PVar(s Series) float64 {
	return variance(s, false)
}

func variance(xs Series, sample bool) float64 {
	if len(xs) <= 1 {
		return math.NaN()
	}

	var m, s float64
	for i, x := range xs {
		mn := m + (x-m)/float64(i+1)
		s += (x - m) * (x - mn)
		m = mn
	}
	if sample {
		return s / float64(len(xs)-1)
	}
	return s / float64(len(xs))
}

// SStd returns the sample standard deviation of the series.
//
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func SStd(s Series) float64 {
	return math.Sqrt(SVar(s))
}

// PStd returns the population standard deviation of the series.
//
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func PStd(s Series) float64 {
	return math.Sqrt(PVar(s))
}

// SSkew returns the sample skew of the series.
//
// NOTE: Not implemented yet.
func SSkew(s Series) float64 {
	panic("not implemented")
}

// PSkew returns the population skew of the series.
//
// NOTE: Not implemented yet.
func PSkew(s Series) float64 {
	panic("not implemented")
}

// SCovar returns the sample covariance of two series s and t.
//
// If the series do not have the same lengths, this function panics.
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func SCovar(s, t Series) float64 {
	n := len(s)
	if n <= 1 {
		return math.NaN()
	}

	u := Mul(Sub1(s, Mean(s)), Sub1(t, Mean(t)))
	return Mean(u) * float64(n) / float64(n-1)
}

// PCovar returns the population covariance of two series s and t.
//
// If the series do not have the same lengths, this function panics.
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func PCovar(s, t Series) float64 {
	n := len(s)
	if n <= 1 {
		return math.NaN()
	}

	return Mean(Mul(s, t)) - Mean(s)*Mean(t)
}

// SCorr returns the sample correlation of two series s and t.
//
// If the series do not have the same lengths, this function panics.
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func SCorr(s, t Series) float64 {
	return SCovar(s, t) / math.Sqrt(SVar(s)*SVar(t))
}

// PCorr returns the poulation correlation of two series s and t.
//
// If the series do not have the same lengths, this function panics.
// If s is empty or has only one element, one cannot speak of variance,
// and NaN is returned.
func PCorr(s, t Series) float64 {
	return PCovar(s, t) / math.Sqrt(PVar(s)*PVar(t))
}

// Autocov returns the sample covariance of s with itself lag values later.
// The series s must be at least 2 longer than lag, else NaN is returned.
func Autocov(s Series, lag int) float64 {
	n := len(s)
	if lag > n-2 {
		return math.NaN()
	}
	return SCovar(s[:n-lag], s[lag:])
}

// Autocorr returns the sample correlation of s with itself lag values later.
// The series s must be at least 2 longer than lag, else NaN is returned.
func Autocorr(s Series, lag int) float64 {
	n := len(s)
	if lag > n-2 {
		return math.NaN()
	}
	return SCorr(s[:n-lag], s[lag:])
}

// Add1 adds f to each value in s and returns a new series.
func Add1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a + f })
}

// Mul1 multiplies f to each value in s and returns a new series.
func Mul1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a * f })
}

// Sub1 subtracts f from each value in s and returns a new series.
func Sub1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a - f })
}

// Div1 divides f from each value in s and returns a new series.
func Div1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a / f })
}

// Add returns the components of s and t added to each other.
func Add(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a + b })
}

// Mul returns the components of s and t multiplied to each other.
func Mul(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a * b })
}

// Sub returns the components of s subtracted by those of t.
func Sub(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a - b })
}

// Div returns the components of s divided by those of t.
func Div(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a / b })
}

// Resize returns s resized to have length n.
//
// If s is empty, the function panics.
// If n is less than the size of s, the first n elements of s is returned.
// If n is greater than the size of s, s is appended to s as often as necessary:
//
// Example:
//
//  Resize([1 2 3], 5) -> [1 2 3 1 2]
//  Resize([1 2 3 4 5], 3) -> [1 2 3]
//
func Resize(s Series, n int) Series {
	m := len(s)
	if n <= m {
		return s[:n].Copy()
	}

	t := make(Series, n)
	for i := 0; i < n; i++ {
		t[i] = s[i%m]
	}
	return t
}

// Fold a series into a single value by repeatedly applying a = f(a, x).
//
// For example, to find the maximum value:
//
//  Fold(s, math.Inf(-1), math.Max)
//
func Fold(s Series, a float64, f func(float64, float64) float64) float64 {
	for _, x := range s {
		a = f(a, x)
	}
	return a
}

// Apply modifies the series by replacing each value v with f(v).
func Apply(s Series, f func(float64) float64) {
	for i, x := range s {
		s[i] = f(x)
	}
}

// Map creates a new series by applying f to each value in s.
func Map(s Series, f func(float64) float64) Series {
	t := make(Series, len(s))
	for i, x := range s {
		t[i] = f(x)
	}
	return t
}

// Map creates a new series by applying f to each value in s and t.
//
// If the series lengths are not the same, the function panics.
func Map2(s, t Series, f func(a, b float64) float64) Series {
	if len(s) != len(t) {
		panic("series lengths must be the same")
	}

	u := make(Series, len(s))
	for i, x := range s {
		u[i] = f(x, t[i])
	}
	return u
}
