package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalc(t *testing.T) {
	// 0.5
	res, err := calc(0.5)
	if assert.NoError(t, err) {
		assert.Equal(t, 18, res)
	}

	// 1.0
	res, err = calc(1)
	if assert.NoError(t, err) {
		assert.Equal(t, 18, res)
	}

	// 2.0
	res, err = calc(2)
	if assert.NoError(t, err) {
		assert.Equal(t, 23, res)
	}

	res, err = calc(20)
	if assert.NoError(t, err) {
		assert.Equal(t, 115, res)
	}
}
