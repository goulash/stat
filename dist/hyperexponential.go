// Copyright (c) 2015, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package dist

import (
	"fmt"
	"math/rand"
)

// HyperExponential distribution with k rates of arrival. {{{
type HyperExponential struct {
	r       *rand.Rand
	lambdas []float64
}

func NewHyperExponential(s rand.Source, lambda []float64) *HyperExponential {
	panic("not implemented")
}

func (e HyperExponential) String() string {
	return fmt.Sprintf("hyper-exponential %v", e.lambdas)
}

func (e HyperExponential) Float64() float64 {
	panic("implement me!")
}

// }}}
