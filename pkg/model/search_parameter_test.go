package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchWithAgeConstraint(t *testing.T) {
	var sp SearchParameter
	SearchWithAgeConstraint(25, 30)(&sp)

	assert.Equal(t, sp.AgeLower, uint(25))
	assert.Equal(t, sp.AgeUpper, uint(30))
}

func TestSearchWithLocationConstraint(t *testing.T) {
	var sp SearchParameter
	SearchWithLocationConstraint(1.2, 9.8, 30)(&sp)

	assert.Equal(t, sp.Latitude, 1.2)
	assert.Equal(t, sp.Longitude, 9.8)
	assert.Equal(t, sp.Radius, uint(30))
}

func TestSearchWithOffset(t *testing.T) {
	var sp SearchParameter
	SearchWithOffset(50)(&sp)

	assert.Equal(t, sp.Offset, uint(50))
}
