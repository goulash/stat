// Copyright (c) 2015, Ben Morgan. All rights reserved.
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

type Series []float64

func (s *Series) Reset()           { *s = make(Series, 0) }
func (s *Series) Append(f float64) { *s = append(*s, f) }

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

func (s Series) Max() float64                         { return Max(s) }
func (s Series) Min() float64                         { return Min(s) }
func (s Series) Mean() float64                        { return Mean(s) }
func (s Series) Median() float64                      { return Median(s) }
func (s Series) SVar() float64                        { return SVar(s) }
func (s Series) PVar() float64                        { return PVar(s) }
func (s Series) SStd() float64                        { return SStd(s) }
func (s Series) PStd() float64                        { return PStd(s) }
func (s Series) SSkew() float64                       { return SSkew(s) }
func (s Series) PSkew() float64                       { return PSkew(s) }
func (s Series) Autocov(lag int) float64              { return Autocov(s, lag) }
func (s Series) Autocorr(lag int) float64             { return Autocorr(s, lag) }
func (s Series) SCovar(t Series) float64              { return SCovar(s, t) }
func (s Series) PCovar(t Series) float64              { return PCovar(s, t) }
func (s Series) SCorr(t Series) float64               { return SCorr(s, t) }
func (s Series) PCorr(t Series) float64               { return PCorr(s, t) }
func (s Series) Apply(f func(float64) float64) Series { return Map(s, f) }
func (s Series) Add1(f float64) Series                { return Add1(s, f) }
func (s Series) Mul1(f float64) Series                { return Mul1(s, f) }
func (s Series) Sub1(f float64) Series                { return Sub1(s, f) }
func (s Series) Div1(f float64) Series                { return Div1(s, f) }
func (s Series) Add(t Series) Series                  { return Add(s, t) }
func (s Series) Mul(t Series) Series                  { return Mul(s, t) }
func (s Series) Sub(t Series) Series                  { return Sub(s, t) }
func (s Series) Div(t Series) Series                  { return Div(s, t) }

func Max(s Series) float64 {
	if len(s) == 0 {
		return 0
	}

	m := s[0]
	for _, x := range s[1:] {
		if m < x {
			m = x
		}
	}
	return m
}

func Min(s Series) float64 {
	if len(s) == 0 {
		return 0
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
func Mean(s Series) float64 {
	if len(s) == 0 {
		return 0
	}

	var m float64
	for i, x := range s {
		m += (x - m) / float64(i+1)
	}
	return m
}

func Median(s Series) float64 {
	if len(s) == 0 {
		return 0
	}

	n := len(s)
	ns := make(Series, n)
	copy(ns, s)
	sort.Float64s(ns)
	if n%2 == 0 {
		return (ns[n/2] + ns[n/2+1]) / 2
	}
	return ns[len(ns)/2]
}

func SVar(s Series) float64 {
	return variance(s, true)
}

func PVar(s Series) float64 {
	return variance(s, false)
}

func variance(xs Series, sample bool) float64 {
	if len(xs) == 0 {
		return 0
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

func SStd(s Series) float64 {
	return math.Sqrt(SVar(s))
}

func PStd(s Series) float64 {
	return math.Sqrt(PVar(s))
}

func SSkew(s Series) float64 {
	panic("not implemented")
}

func PSkew(s Series) float64 {
	panic("not implemented")
}

func SCovar(s, t Series) float64 {
	if len(s) != len(t) {
		panic("series not equal length")
	}

	n := float64(len(s))
	u := Mul(Sub1(s, Mean(s)), Sub1(t, Mean(t)))
	return Mean(u) * n / (n - 1)
}

// PCovar returns the population covariance of two series s and t.
func PCovar(s, t Series) float64 {
	return Mean(Mul(s, t)) - Mean(s)*Mean(t)
}

func SCorr(s, t Series) float64 {
	return SCovar(s, t) / math.Sqrt(SVar(s)*SVar(t))
}

func PCorr(s, t Series) float64 {
	return PCovar(s, t) / math.Sqrt(PVar(s)*PVar(t))
}

func Autocov(s Series, lag int) float64 {
	if lag >= len(s) {
		panic("lag is too long")
	}
	return SCovar(s[:len(s)-lag], s[lag:])
}

func Autocorr(s Series, lag int) float64 {
	if lag >= len(s) {
		panic("lag is too long")
	}
	return SCorr(s[:len(s)-lag], s[lag:])
}

func Tail(s Series, n int) Series {
	panic("implement me!")
}

func Head(s Series, n int) Series {
	panic("implemetn me!")
}

func Add1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a + f })
}

func Mul1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a * f })
}

func Sub1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a - f })
}

func Div1(s Series, f float64) Series {
	return Map(s, func(a float64) float64 { return a / f })
}

func Add(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a + b })
}

func Mul(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a * b })
}

func Sub(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a - b })
}

func Div(s, t Series) Series {
	return Map2(s, t, func(a, b float64) float64 { return a / b })
}

func Map(s Series, f func(float64) float64) Series {
	t := make(Series, len(s))
	for i, x := range s {
		t[i] = f(x)
	}
	return t
}

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
