package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

const driveMapperPath = "/dev/mapper/%s"
const driveMntPath = "/mnt/%s"
const mountsPath = "/proc/mounts"

// DriveStatus contains information about specific drive
type DriveStatus struct {
	Name      string `json:"name"`      // Human redable name
	UUID      string `json:"uuid"`      // UUID from /dev/disks/by-uuid
	Encrypted bool   `json:"encrypted"` // True if the drive is configured as encrypted
	Mounted   bool   `json:"mounted"`   // True if mounted
	Opened    bool   `json:"opened"`    // True in case the drive is encrypted and cryptsetup open finished without any error
	Smart     string `json:"smart"`     // String with: ok (nothing wrong is happening), warning (good to look for a new one), error (replace the drive ASAP), Unknown (we don't know)
}

// IsReady returns true if the drive is mounted and opened
func (d *DriveStatus) IsReady() bool {
	return d.Mounted && d.Opened
}

// GetDriveStatus returns filled DriveStatus struct where information about the drive is saved
func GetDriveStatus(driveIndex uint) (*DriveStatus, error) {
	var status DriveStatus

	drive := config.Drives[driveIndex]

	status.Name = drive.Name
	status.UUID = drive.UUID
	status.Encrypted = drive.Encrypted
	status.Smart = "unknown"

	// Check if the drive is openned
	if _, err := os.Stat(fmt.Sprintf(driveMapperPath, status.Name)); err == nil && status.Encrypted {
		status.Opened = true
	} else {
		status.Opened = false
	}

	// Check if the drive is mounted
	content, err := ioutil.ReadFile(mountsPath)
	if err != nil {
		return &status, err
	}
	mounts := bytes.Split(content, []byte("\n"))

	for i := range mounts {
		src := bytes.Split(mounts[i], []byte(" "))[0]
		if string(src) == fmt.Sprintf(driveMapperPath, status.Name) {
			status.Mounted = true
			break
		}
	}

	return &status, nil
}

// MountDrive mounts drive defined in status variables
func (d *DriveStatus) MountDrive() ([]byte, error) {
	src := fmt.Sprintf(driveMapperPath, d.Name)
	target := fmt.Sprintf(driveMntPath, d.Name)
	args := []string{src, target}
	out, err := runCommand("/bin/mount", args, nil)
	return out, err
}

// UmountDrive umounts drive defined in status variables
func (d *DriveStatus) UmountDrive() ([]byte, error) {
	target := fmt.Sprintf(driveMntPath, d.Name)
	args := []string{target}

	out, err := runCommand("/bin/umount", args, nil)
	return out, err
}

// OpenDrive opens encrypted device
func (d *DriveStatus) OpenDrive(password string) ([]byte, error) {
	args := []string{"open", "/dev/disk/by-uuid/" + d.UUID, d.Name}
	out, err := runCommand("/sbin/cryptsetup", args, []byte(password+"\n"))
	return out, err
}

// CloseDrive closes encrypted device
func (d *DriveStatus) CloseDrive() ([]byte, error) {
	src := fmt.Sprintf(driveMapperPath, d.Name)
	args := []string{"close", src}

	out, err := runCommand("/sbin/cryptsetup", args, nil)
	return out, err
}
