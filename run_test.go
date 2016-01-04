// Copyright (c) 2016, Ben Morgan. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package stat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSeries(z *testing.T) {
	assert := assert.New(z)

	tests := []Series{
		{1, 2, 3, 4, 5},
		{0.38809179, 0.94113008, 0.15350705, 0.03311646, 0.68168087, 0.21719990},
		{0.32123922, 0.57085251, 0.53576882, 0.38965630, 0.27487263, 0.90783122},
		{0, 0, 0, 0, 0},
		{-1, -2, -3},
	}

	for _, t := range tests {
		var r Run
		for _, f := range t {
			r.Add(f)
		}

		assert.Equal(t.Min(), r.Min())
		assert.Equal(t.Max(), r.Max())
		assert.Equal(t.Mean(), r.Mean())
		assert.Equal(t.Var(), r.Var())
		assert.Equal(t.Std(), r.Std())
		assert.Equal(t.VarP(), r.VarP())
		assert.Equal(t.StdP(), r.StdP())
	}
}
