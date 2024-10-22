package main

import (
	_ "embed"
	"io"
	"sync"
	"text/template"
)

//go:embed templates/default.go.tpl
var defaultTemplate string

type Generator struct {
	Template string

	loadOnce sync.Once
	loadErr  error
	template *template.Template
}

func (g *Generator) Generate(w io.Writer, file *File) error {
	g.loadOnce.Do(func() {
		tpl := g.Template
		if tpl == "" {
			tpl = defaultTemplate
		}
		g.template, g.loadErr = template.New("tpl").Parse(tpl)
	})
	if g.loadErr != nil {
		return g.loadErr
	}

	if err := g.template.Execute(w, file); err != nil {
		return err
	}
	return nil
}
