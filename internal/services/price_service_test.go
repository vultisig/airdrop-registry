package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCMCMap(t *testing.T) {
	pr, err := NewPriceResolver()
	assert.Nil(t, err)
	assert.NotNil(t, pr)
}
