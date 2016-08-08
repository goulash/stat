// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

type Null struct{}

func NewNull() *Null            { return &Null{} }
func (n Null) Int63() int64     { return 0 }
func (n Null) Float64() float64 { return 0.0 }
func (n Null) Mean() float64    { return 0.0 }
func (n Null) Var() float64     { return 0.0 }
func (n Null) Std() float64     { return 0.0 }
func (n Null) String() string   { return "null" }
