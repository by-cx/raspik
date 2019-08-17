package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", "")
}

func getDrives(c echo.Context) error {
	notReadyFilter := false
	if c.QueryParam("not_ready") == "1" {
		notReadyFilter = true
	}

	driveStatuses := make([]DriveStatus, 0, len(config.Drives))

	for drive := range config.Drives {
		status, err := GetDriveStatus(uint(drive))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, err.Error())
		}

		if (notReadyFilter && !status.IsReady()) || !notReadyFilter {
			driveStatuses = append(driveStatuses, *status)
		}
	}

	return c.JSONPretty(http.StatusOK, driveStatuses, "    ")
}

func postDrivesUnlock(c echo.Context) error {
	var passwords []string

	err := c.Bind(&passwords)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Load drives status
	driveStatuses := make([]DriveStatus, 0, len(config.Drives))

	for idx, drive := range config.Drives {
		status, err := GetDriveStatus(uint(idx))
		if err != nil {
			return echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("Error while getting info about '%s' | '%s'", drive.Name, err.Error()))
		}
		if !status.IsReady() {
			driveStatuses = append(driveStatuses, *status)
		}
	}

	if len(driveStatuses) != len(passwords) {
		return echo.NewHTTPError(http.StatusBadRequest, "Number of password doesn't match number of drives")
	}

	for idx, drive := range driveStatuses {

		if drive.Encrypted && !drive.Opened {
			_, err := drive.OpenDrive(passwords[idx])
			if err != nil {
				return echo.NewHTTPError(
					http.StatusBadGateway,
					fmt.Sprintf("Error while opening '%s' | '%s' | Is your password correct?", drive.Name, err.Error()),
				)
			}
		}
		if !drive.Mounted {
			out, err := drive.MountDrive()
			if err != nil {
				return echo.NewHTTPError(
					http.StatusBadGateway,
					fmt.Sprintf("Error while mounting '%s' | '%s' | '%s'", drive.Name, err.Error(), string(out)),
				)
			}
		}
	}

	// If ready then run all enabled service
	ready := true
	for idx, drive := range config.Drives {
		status, err := GetDriveStatus(uint(idx))
		if err != nil {
			return echo.NewHTTPError(
				http.StatusBadGateway,
				fmt.Sprintf("Error while getting info about '%s' | '%s'", drive.Name, err.Error()),
			)
		}
		if !status.IsReady() {
			ready = false
			break
		}
	}

	if ready {
		for _, user := range config.Users {
			if user.Services.Syncthing.Enabled {
				if err := startService("syncing@" + user.Name); err != nil {
					return echo.NewHTTPError(
						http.StatusBadGateway,
						fmt.Sprintf("Error while starting syncthing | '%s'", err.Error()),
					)
				}
			}
		}
	}

	return c.JSONPretty(http.StatusOK, map[string]string{"message": "ok"}, "    ")
}

func shares(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, config.Shares, "    ")
}

func users(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, config.Users, "    ")
}

func general(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, config.General, "    ")
}

func backup(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, config.Backup, "    ")
}
