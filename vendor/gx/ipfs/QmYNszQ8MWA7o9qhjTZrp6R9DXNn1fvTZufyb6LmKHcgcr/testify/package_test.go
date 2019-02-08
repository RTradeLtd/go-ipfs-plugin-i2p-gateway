package testify

import (
	"gx/ipfs/QmYNszQ8MWA7o9qhjTZrp6R9DXNn1fvTZufyb6LmKHcgcr/testify/assert"
	"testing"
)

func TestImports(t *testing.T) {
	if assert.Equal(t, 1, 1) != true {
		t.Error("Something is wrong.")
	}
}
