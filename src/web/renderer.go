package main

import (
	"io"
	"text/template"

	"github.com/labstack/echo"
)

// Template renderer
type Template struct {
	templates *template.Template
}

// Render connects data with the template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
