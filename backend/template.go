package backend

import (
	"fmt"
	"github.com/kohirens/www/storage"
	"io"
	"maps"
	"os"
	"text/template"
)

type Renderer struct {
	store    storage.Storage
	location string
	suffix   string
	Vars     map[string]string
}

type TemplateManager interface {
	// AddVar Add an item to the variable map.
	AddVar(k, v string)
	// AppendVars Appends a map to the Renderer.Vars map.
	//
	//	NOTE: When a key matches an existing key it will overwrite its value.
	AppendVars(vars map[string]string)
	// Load A template, but it will not render it, instead, the template.Template
	// object is returned so that you can render it when you want.
	Load(name string) (*template.Template, error)
	// Render Write a templates' content to a writer. You can provide vars
	// as a type `map[string]string` of key-value pairs; which will be used to fill
	// in string placeholders. Nothing more complex is supported at this time.
	// Also, remember that maps are by default passed by reference, so there is
	// no need to pass vars as a pointer.
	Render(name string, w io.Writer, vars map[string]string) error
}

type Variables map[string]string

const ps = string(os.PathSeparator)

func NewTemplateManager(store storage.Storage, location, suffix string) TemplateManager {
	return &Renderer{
		location: location,
		store:    store,
		suffix:   suffix,
		Vars:     make(map[string]string),
	}
}

// AddVar Add an item to the variable map.
func (m *Renderer) AddVar(k, v string) {
	m.Vars[k] = v
}

// AppendVars Appends a map to the Renderer.Vars map.
//
//	NOTE: When a key matches an existing key it will overwrite its value.
func (m *Renderer) AppendVars(vars map[string]string) {
	maps.Copy(m.Vars, vars)
}

// Load A template, but it will not render it, instead, the template.Template
// object is returned so that you can render it when you want.
func (m *Renderer) Load(name string) (*template.Template, error) {
	filename := m.location + name + "." + m.suffix
	if len(name) > 0 && name[0] != '/' {
		filename = m.location + ps + name + "." + m.suffix
	}

	Log.Infof("load template %v.", filename)

	tmplContent, e1 := m.store.Load(filename)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.FileNotFound, filename, e1.Error())
	}

	t, e2 := template.New(filename).Parse(string(tmplContent))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.TemplateParse, e2.Error())
	}

	Log.Dbugf(stdout.TemplateLoad, filename)

	return t, nil
}

// Render Write a templates' content to a writer. You can provide vars
// as a type `map[string]string` of key-value pairs; which will be used to fill
// in string placeholders. Nothing more complex is supported at this time.
// Also, remember that maps are by default passed by reference, so there is
// no need to pass vars as a pointer.
func (m *Renderer) Render(name string, w io.Writer, vars map[string]string) error {
	t, e1 := m.Load(name)
	if e1 != nil {
		return e1
	}

	// Combine vars with any previously added to the manager.
	if vars != nil {
		m.AppendVars(vars)
	}

	return t.Execute(w, m.Vars)
}
