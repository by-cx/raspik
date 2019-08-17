package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/yaml.v2"
)

// ConfigPath sets where the config file is saved
//const ConfigPath = "/etc/raspirack/config.yml"
const ConfigPath = "config.dev.yml"

var config *Config

func init() {
	config = loadConfig()
}

func loadConfig() *Config {
	var config Config

	data, err := ioutil.ReadFile(ConfigPath)
	check(err)

	err = yaml.Unmarshal(data, &config)
	check(err)

	return &config
}

func printConfig() {
	output, err := yaml.Marshal(config)
	check(err)

	fmt.Printf("--- config dump:\n%s\n\n", string(output))
}

// Template renderer
type Template struct {
	templates *template.Template
}

// Render connects data with the template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Use(middleware.Logger())

	e.GET("/", indexHandler)
	e.GET("/drives", getDrivesHandler)
	e.POST("/drives/unlock", postDrivesUnlockHandler)
	e.GET("/sharesHandler", sharesHandler)
	e.GET("/usersHandler", usersHandler)
	e.GET("/generalHandler", generalHandler)
	e.GET("/backupHandler", backupHandler)

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	e.Logger.Fatal(e.Start(":1323"))
}
