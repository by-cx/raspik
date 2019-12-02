package main

import (
	"net/http"
	"text/template"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("src/templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.Logger())

	e.Static("/static", "src/static/")
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", "users")
	})
	e.GET("/config", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", "config")
	})
	e.GET("/unlock", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", "unlock")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
