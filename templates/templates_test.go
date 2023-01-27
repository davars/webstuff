package templates

import (
	"bytes"
	"embed"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed tests/valid/templates/*
var ValidTemplates embed.FS

func TestNewCollection(t *testing.T) {
	validFS, err := fs.Sub(ValidTemplates, "tests/valid")
	assert.NoError(t, err)
	coll, err := NewCollection(validFS, nil)
	assert.NoError(t, err)
	var buf bytes.Buffer
	assert.NoError(t, coll.Render(&buf, "page1.gohtml", nil))
	contents := buf.String()
	assert.Contains(t, contents, "Title 1")
	assert.Contains(t, contents, "Main 1")
	assert.Contains(t, contents, "This is the nav")
	assert.Contains(t, contents, "I'm the base.")
}
