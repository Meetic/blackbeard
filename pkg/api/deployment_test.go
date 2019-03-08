package api_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKill(t *testing.T) {
	errs := blackbeard.Kill("foo", []string{"test"})
	assert.Empty(t, errs)
}

func TestKillKO(t *testing.T) {
	errs := blackbeard.Kill("foo", []string{"bar"})
	assert.NotEmpty(t, errs)
}
