package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootIDIsZero(t *testing.T) {
	assert.Equal(t, 0, RootID)
}

func TestRootItemHasRootID(t *testing.T) {
	assert.Equal(t, int64(RootID), RootItem.ID)
}
