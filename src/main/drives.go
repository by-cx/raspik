package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

const driveMapperPath = "/dev/mapper/%s"
const driveMntPath = "/mnt/%s"
const mountsPath = "/proc/mounts"

// DriveStatus contains information about specific drive
type DriveStatus struct {
	Name      string // Human redable name
	UUID      string // UUID from /dev/disks/by-uuid
	Password  string // Password to dencrypt the drive
	Encrypted bool   // True if the drive is configured as encrypted
	Mounted   bool   // True if mounted
	Opened    bool   // True in case the drive is encrypted and cryptsetup open finished without any error
	Smart     string // String with: ok (nothing wrong is happening), warning (good to look for a new one), error (replace the drive ASAP), Unknown (we don't know)
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
func MountDrive(status *DriveStatus) ([]byte, error) {
	src := fmt.Sprintf(driveMapperPath, status.Name)
	target := fmt.Sprintf(driveMntPath, status.Name)
	cmd := fmt.Sprintf("mount -o default %s %s", src, target)

	out, err := exec.Command(cmd).Output()
	return out, err
}

// UmountDrive umounts drive defined in status variables
func UmountDrive(status *DriveStatus) ([]byte, error) {
	target := fmt.Sprintf(driveMntPath, status.Name)
	cmd := fmt.Sprintf("umount %s", target)

	out, err := exec.Command(cmd).Output()
	return out, err
}

// OpenDrive opens encrypted device
func openDrive(status *DriveStatus) ([]byte, error) {
	cmd := fmt.Sprintf("cryptsetup open /dev/disk/by-uuid/%s %s", status.Name, status.Name)
	subprocess := exec.Command(cmd)

	writer, err := subprocess.StdinPipe()
	if err != nil {
		return []byte(""), err
	}

	_, err = writer.Write([]byte(status.Password + "\n"))
	if err != nil {
		return []byte(""), err
	}

	out, err := subprocess.Output()
	return out, err
}

// CloseDrive closes encrypted device
func CloseDrive(status *DriveStatus) ([]byte, error) {
	src := fmt.Sprintf(driveMapperPath, status.Name)
	cmd := fmt.Sprintf("cryptsetup close %s", src)

	out, err := exec.Command(cmd).Output()
	return out, err
}
