package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/yaml.v2"
)

// ConfigPath sets where the config file is saved
const ConfigPath = "/etc/raspirack/config.yml"

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
		templates: template.Must(template.ParseGlob("src/templates/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
	}))
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", "")
	})
	e.GET("/device/status", func(c echo.Context) error {
		rpi := RPi{}
		cpuTemperature, err := rpi.CPUTemperature()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.JSONPretty(http.StatusInternalServerError, DeviceStatus{
			CPUTemperature: cpuTemperature,
		}, jsonIdent)
	})
	e.GET("/drives", func(c echo.Context) error {
		notReadyFilter := false
		if c.QueryParam("not_ready") == "1" {
			notReadyFilter = true
		}

		var driveStatuses []DriveStatus = []DriveStatus{}

		for drive := range config.Drives {
			status, err := GetDriveStatus(uint(drive))
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			if (notReadyFilter && !status.IsReady()) || !notReadyFilter {
				driveStatuses = append(driveStatuses, *status)
			}
		}

		return c.JSONPretty(http.StatusOK, driveStatuses, jsonIdent)
	})
	e.POST("/drives/unlock", func(c echo.Context) error {
		passwordRaw, err := ioutil.ReadAll(c.Request().Body)
		if err != nil && err.Error() != "EOF" {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		password := string(passwordRaw)

		// Load drives status
		var driveStatuses = []DriveStatus{}

		// Read status of all drives
		for idx, drive := range config.Drives {
			status, err := GetDriveStatus(uint(idx))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, "Error while getting info about "+drive.Name+" | "+err.Error())
			}
			if !status.IsReady() {
				driveStatuses = append(driveStatuses, *status)
			}
		}

		// Open the ones that are not openned and mount them
		for _, drive := range driveStatuses {
			if drive.Encrypted && !drive.Opened {
				_, err := drive.OpenDrive(password)
				if err != nil {
					return echo.NewHTTPError(http.StatusBadGateway, "Error while opening "+drive.Name+" | "+err.Error()+" | Is your password correct?")
				}
			}
			if !drive.Mounted {
				out, err := drive.MountDrive()
				if err != nil {
					return echo.NewHTTPError(http.StatusBadGateway, "Error while mounting "+drive.Name+" | "+err.Error()+" | "+string(out))
				}
			}
		}

		// If ready then run all enabled service
		ready := true
		for idx, drive := range config.Drives {
			status, err := GetDriveStatus(uint(idx))
			if err != nil {
				return echo.NewHTTPError(http.StatusBadGateway, "Error while getting info about "+drive.Name+" | "+err.Error())
			}
			if !status.IsReady() {
				ready = false
				break
			}
		}

		// Start Syncthing
		if ready {
			for _, user := range config.Users {
				if user.Services.Syncthing.Enabled {
					err := startService("syncthing@" + user.Name)
					if err != nil {
						return echo.NewHTTPError(http.StatusBadGateway, "Error while starting syncthing | "+err.Error())
					}
				}
			}
		}

		return c.JSONPretty(http.StatusOK, map[string]string{"message": "ok"}, "    ")
	})
	e.GET("/shares", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, config.Shares, "    ")
	})
	e.GET("/users", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, config.Users, "    ")
	})
	e.GET("/general", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, config.General, "    ")
	})
	e.GET("/backup", func(c echo.Context) error {
		return c.JSONPretty(http.StatusOK, config.Backup, "    ")
	})

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	e.Logger.Fatal(e.Start(":1323"))
}
