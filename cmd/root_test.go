package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAskForConfirmation(t *testing.T) {
	retNo := askForConfirmation("test", strings.NewReader("no\n"))
	assert.False(t, retNo)
	retYes := askForConfirmation("test", strings.NewReader("yes\n"))
	assert.True(t, retYes)
}
