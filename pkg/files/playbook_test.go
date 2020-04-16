package files_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Meetic/blackbeard/pkg/files"
)

func TestGetTemplate(t *testing.T) {
	r := files.NewPlaybookRepository("templates_test", "templates_test")

	tpls, err := r.GetTemplate()

	assert.Len(t, tpls, 2)
	assert.Equal(t, "01_template.yml", tpls[0].Name)
	assert.Equal(t, "02_template.yml", tpls[1].Name)
	assert.Nil(t, err)
}

func TestGetTemplateNotFound(t *testing.T) {
	r := files.NewPlaybookRepository(".", ".")

	tpls, err := r.GetTemplate()

	assert.Len(t, tpls, 0)
	assert.NotNil(t, err)
}
