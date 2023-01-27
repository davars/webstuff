package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"path/filepath"
)

// A Collection is a set of pages that share a base layout, partials, and template funcs
type Collection struct {
	ts map[string]*template.Template
}

// NewCollection takes a fs.FS containing the template sources.  The filesystem is expected
// to be in the following format:
//
//	templates/
//	├─ pages/
//	│  ├─ page1.gohtml
//	│  ├─ page2.gohtml
//	├─ partials/
//	│  ├─ nav.gohtml
//	│  ├─ footer.gohtml
//	├─ base.gohtml
//
// base.gohtml is a required filename and defines the top-level template.  It refers to content
// blocks defined by the partials and pages. The filenames for the partials don't matter, only
// the names of the templates they define. The pages are referred by their filename in the Render
// method.
func NewCollection(filesystem fs.FS, funcMap template.FuncMap) (Collection, error) {
	set := Collection{
		ts: map[string]*template.Template{},
	}

	if _, err := filesystem.Open("templates/base.gohtml"); err != nil {
		return Collection{}, fmt.Errorf("base.gohtml not found in filesystem")
	}

	pages, err := fs.Glob(filesystem, "templates/pages/*.gohtml")
	if err != nil {
		return Collection{}, err
	}
	if len(pages) == 0 {
		return Collection{}, fmt.Errorf("no pages found")
	}
	partials, err := fs.Glob(filesystem, "templates/partials/*.gohtml")
	if err != nil {
		return Collection{}, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts := append([]string{"templates/base.gohtml"}, append(partials, page)...)
		t, err := template.New("").Funcs(funcMap).ParseFS(filesystem, ts...)
		if err != nil {
			return Collection{}, err
		}

		set.ts[name] = t
	}
	return set, nil
}

func (tc Collection) Render(w io.Writer, page string, data any) error {
	if tc.ts == nil {
		panic(fmt.Errorf("must use NewCollection to create a Collection"))
	}
	t, ok := tc.ts[page]
	if !ok {
		return fmt.Errorf("the template %s does not exist", page)
	}

	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	return err
}
