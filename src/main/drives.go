package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

const driveMapperPath = "/dev/mapper/%s"
const driveMntPath = "/mnt/%s"
const mountsPath = "/proc/mounts"

// DriveStatus contains information about specific drive
type DriveStatus struct {
	Name             string             `json:"name"`      // Human redable name
	UUID             string             `json:"uuid"`      // UUID from /dev/disks/by-uuid
	Encrypted        bool               `json:"encrypted"` // True if the drive is configured as encrypted
	Mounted          bool               `json:"mounted"`   // True if mounted
	Opened           bool               `json:"opened"`    // True in case the drive is encrypted and cryptsetup open finished without any error
	Health           string             `json:"health"`    // String with: ok (nothing wrong is happening), warning (good to look for a new one), error (replace the drive ASAP), Unknown (we don't know)
	Failures         []string           `json:"failures"`  // Health status translated into failure messages, if empty, all is good
	Raid             bool               `json:"raid"`
	RaidDevices      []string           `json:"raid_devices"` // UUID of block devices where the Btrfs RAID is
	HealthRAW        []Health           `json:"health_raw"`
	UsedBytes        uint64             `json:"used_bytes"`
	TotalBytes       uint64             `json:"total_bytes"`
	ScrubRunning     bool               `json:"scrub_running"`
	ScrubErrors      int64              `json:"scrub_errors"`
	FilesystemErrors []FileSystemErrors `json:"filesystem_errors"`
}

// IsReady returns true if the drive is mounted and opened
func (d *DriveStatus) IsReady() bool {
	return d.Mounted && d.Opened
}

// GetTarget returns path where the device is mounted
func (d *DriveStatus) GetTarget() string {
	return fmt.Sprintf(driveMntPath, d.Name)
}

// GetUUID returns UUIDs of all underlaying drives
func (d *DriveStatus) GetUUID() []string {
	uuids := []string{d.UUID}
	if d.Raid {
		uuids = d.RaidDevices
	}
	return uuids
}

func (d *DriveStatus) getRootBlockDevices() ([]string, error) {
	uuids := d.GetUUID()

	var devices []string

	for _, blockDeviceUUID := range uuids {
		pointedTo, err := os.Readlink("/dev/disk/by-uuid/" + blockDeviceUUID)
		if err != nil {
			return []string{}, err
		}
		parts := strings.Split(pointedTo, "/")
		if len(parts) > 0 {
			devices = append(devices, parts[len(parts)-1])
		} else {
			return []string{}, errors.New("symlink cannot be resolved")
		}

	}

	return devices, nil
}

// ReadHealth reads S.M.A.R.T. health status
func (d *DriveStatus) ReadHealth() ([]Health, error) {
	var healths []Health

	devices, err := d.getRootBlockDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if len(device) <= 3 {
			return nil, errors.New("strange format of device name (" + device + ")")
		}

		device = device[0:3]

		health := Health{
			Device:           device,
			RelocatedSectors: -1,
			PendingSectors:   -1,
			Temperature:      -99,
		}

		output, err := runCommand("/usr/sbin/smartctl", []string{"-a", "/dev/" + device}, []byte(""))
		if err != nil {
			return nil, err
		}

		// Check pending and relocated sectors and temperature
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Current_Pending_Sector") {
				fields := strings.Fields(line)
				if len(fields) >= 10 {
					value, err := strconv.Atoi(fields[9])
					if err != nil {
						return nil, err
					}
					health.PendingSectors = value
				}
			} else if strings.Contains(line, "Reallocated_Sector_Ct") {
				fields := strings.Fields(line)
				if len(fields) >= 10 {
					value, err := strconv.Atoi(fields[9])
					if err != nil {
						return nil, err
					}
					health.RelocatedSectors = value
				}
			} else if strings.Contains(line, "Temperature_Celsius") {
				fields := strings.Fields(line)
				if len(fields) >= 10 {
					value, err := strconv.Atoi(fields[9])
					if err != nil {
						return nil, err
					}
					health.Temperature = value
				}
			}
		}

		// Check health status in S.M.A.R.T.
		output, err = runCommand("/usr/sbin/smartctl", []string{"-H", "/dev/" + device}, []byte(""))
		if err != nil {
			return nil, err
		}

		lines = strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "SMART overall-health self-assessment test result") {
				fields := strings.Fields(line)
				if fields[len(fields)-1] == "PASSED" {
					health.SMARTHealth = true
				}
			}
		}

		healths = append(healths, health)
	}

	return healths, nil
}

// ReadFsErrors returns number of filesystem errors.
// The result is sum of errors on all devices.
func (d *DriveStatus) ReadFsErrors() ([]FileSystemErrors, error) {
	var fileSystemErrors []FileSystemErrors

	output, err := runCommand("/bin/btrfs", []string{"device", "stats", d.GetTarget()}, []byte(""))
	if err != nil {
		return fileSystemErrors, nil
	}

	// Source is device in mapper if the data is encrypted, otherwise it's just path to the block device
	// It's also first field in output of btrfs device stats.
	source := ""
	var fsErrors FileSystemErrors
	output = bytes.TrimSpace(output)
	lines := strings.Split(string(output), "\n")
	for idx, line := range lines {
		fields := strings.Fields(line)

		if len(fields) != 2 {
			return fileSystemErrors, errors.New("unexpected output of btrfs device stats")
		}

		parsedSource := strings.Split(fields[0], ".")

		if len(parsedSource) != 2 {
			return fileSystemErrors, errors.New("unexpected output of btrfs device stats")
		}

		if source == "" {
			source = parsedSource[0]
		}

		value, err := strconv.Atoi(fields[1])
		if err != nil {
			return fileSystemErrors, err
		}

		if strings.Contains(line, "write_io_errs") {
			fsErrors.WriteIOErrors += value
		} else if strings.Contains(line, "read_io_errs") {
			fsErrors.ReadIOErrors += value
		} else if strings.Contains(line, "flush_io_errs") {
			fsErrors.FlushIOErrors += value
		} else if strings.Contains(line, "corruption_errs") {
			fsErrors.CorruptionErrors += value
		} else if strings.Contains(line, "generation_errs") {
			fsErrors.GenerationErrsor += value
		}

		if parsedSource[0] != source || len(lines) == idx+1 {
			fileSystemErrors = append(fileSystemErrors, fsErrors)

			fsErrors.WriteIOErrors = 0
			fsErrors.ReadIOErrors = 0
			fsErrors.FlushIOErrors = 0
			fsErrors.CorruptionErrors = 0
			fsErrors.GenerationErrsor = 0

			source = parsedSource[0]
		}
	}

	return fileSystemErrors, nil

}

// Df returns occupied bytes and total bytes available on this drive or error in case of trouble.
func (d *DriveStatus) Df() (uint64, uint64, error) {
	if d.Mounted {
		stats, err := disk.Usage(d.GetTarget())
		if err != nil {
			return 0, 0, err
		}

		return stats.Used, stats.Total, nil
	}

	return 0, 0, nil
}

// ScrubData returns number of errors found during last round of scrubbing. Returned values
// are true or false for running, number of errors and error if something goes wrong.
func (d *DriveStatus) ScrubData() (bool, int64, error) {
	// Running
	// scrub status for 751e6d40-a421-439b-bb73-ec9a53dc100e
	// scrub started at Sat Oct 17 14:54:57 2020, running for 02:31:16
	// total bytes scrubbed: 1.17TiB with 0 errors

	// Done
	// scrub status for 751e6d40-a421-439b-bb73-ec9a53dc100e
	// scrub started at Sat Oct 17 14:54:57 2020 and finished after 09:10:31
	// total bytes scrubbed: 4.29TiB with 0 errors

	running := false
	errors := int64(0)

	if !d.Mounted {
		return running, -1, nil
	}

	output, err := runCommand("/bin/btrfs", []string{"scrub", "status", d.GetTarget()}, []byte(""))
	if err != nil {
		return running, errors, err
	}

	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "running for") {
			running = true
		}
		if strings.Contains(line, "total bytes scrubbed") {
			fields := strings.Fields(line)
			value, err := strconv.Atoi(fields[len(fields)-2])
			if err != nil {
				return running, errors, err
			}
			errors = int64(value)
		}
	}

	return running, errors, nil

}

// MountDrive mounts drive defined in status variables
func (d *DriveStatus) MountDrive() ([]byte, error) {
	log.Println("Mounting drive " + d.Name)
	name := d.Name

	if d.Raid {
		name += "_0" // We need only first btrfs device to mount all of them
	}

	src := fmt.Sprintf(driveMapperPath, name)
	target := d.GetTarget()
	args := []string{src, target}
	out, err := runCommand("/bin/mount", args, nil)
	return out, err
}

// UmountDrive umounts drive defined in status variables
func (d *DriveStatus) UmountDrive() ([]byte, error) {
	log.Println("Umounting drive " + d.Name)

	name := d.Name

	if d.Raid {
		name += "0" // We need only first btrfs device to mount all of them
	}

	target := d.GetTarget()
	args := []string{target}

	out, err := runCommand("/bin/umount", args, nil)
	return out, err
}

// OpenDrive opens encrypted device
func (d *DriveStatus) OpenDrive(password string) ([]byte, error) {
	// If there is a raid, we have to loop over all backend devices and call cryptsetup open on each of them
	if d.Raid {
		var fullOut []byte

		if len(d.RaidDevices) == 0 {
			return []byte(""), errors.New("no raid devices configured but raid enabled")
		}

		for idx, backendDeviceUUID := range d.RaidDevices {
			log.Println("Opening encrypted drive " + d.Name + "_" + strconv.Itoa(idx))
			args := []string{"open", "/dev/disk/by-uuid/" + backendDeviceUUID, d.Name + "_" + strconv.Itoa(idx)}
			out, err := runCommand("/sbin/cryptsetup", args, []byte(password+"\n"))
			if err != nil {
				return out, err
			}

			fullOut = append(fullOut, out...)
			fullOut = append(fullOut, '\n')
		}
		return fullOut, nil
	} else {
		log.Println("Opening encrypted drive " + d.Name)
		args := []string{"open", "/dev/disk/by-uuid/" + d.UUID, d.Name}
		out, err := runCommand("/sbin/cryptsetup", args, []byte(password+"\n"))
		return out, err
	}
}

// CloseDrive closes encrypted device
func (d *DriveStatus) CloseDrive() ([]byte, error) {
	// If there is a raid, we have to loop over all backend devices and call cryptsetup close on each of them
	if d.Raid {
		var fullOut []byte
		for idx := range d.RaidDevices {
			log.Println("Closing encrypted drive " + d.Name + "_" + strconv.Itoa(idx))
			src := fmt.Sprintf(driveMapperPath, d.Name+"_"+strconv.Itoa(idx))
			args := []string{"close", src}

			out, err := runCommand("/sbin/cryptsetup", args, nil)
			if err != nil {
				return out, err
			}

			fullOut = append(fullOut, out...)
			fullOut = append(fullOut, '\n')
		}
		return fullOut, nil
	} else {
		log.Println("Closing encrypted drive " + d.Name)
		src := fmt.Sprintf(driveMapperPath, d.Name)
		args := []string{"close", src}

		out, err := runCommand("/sbin/cryptsetup", args, nil)
		return out, err
	}

}

// GetDriveStatus returns filled DriveStatus struct where information about the drive is saved.
func GetDriveStatus(driveIndex uint) (*DriveStatus, error) {
	var status DriveStatus

	drive := config.Drives[driveIndex]

	status.Name = drive.Name
	status.UUID = drive.UUID
	status.Encrypted = drive.Encrypted
	status.Health = "unknown"
	status.Raid = drive.Raid
	status.RaidDevices = drive.RaidDevices

	srcDriveName := status.Name
	if status.Raid {
		// In case of raid we check if the first device is opened or not.
		// TODO: We should check all of them.
		srcDriveName += "_0"
	}

	// Check if the drive is openned
	if _, err := os.Stat(fmt.Sprintf(driveMapperPath, srcDriveName)); err == nil && status.Encrypted {
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
		if string(src) == fmt.Sprintf(driveMapperPath, srcDriveName) {
			status.Mounted = true
			break
		}
	}

	// Filesystem errors
	status.FilesystemErrors, err = status.ReadFsErrors()
	if err != nil {
		return nil, err
	}

	// The checks at the beggining are more important compare to the checks at the end.
	// We don't care about temperature when there are bad sectors on one of the drives
	// for example.

	// S.M.A.R..T. health
	status.HealthRAW, err = status.ReadHealth()
	if err != nil {
		return nil, err
	}

	failures := []string{}
	for _, healthRAW := range status.HealthRAW {
		if healthRAW.PendingSectors > 0 && healthRAW.RelocatedSectors > 0 {
			failures = append(failures, "bad sectors detected")
		}
		if !healthRAW.SMARTHealth {
			failures = append(failures, "S.M.A.R.T. self-test failure")
		}
	}

	// Check for file system errors
	for _, fsError := range status.FilesystemErrors {
		if fsError.Total() > 0 {
			failures = append(failures, "filesystem errors detected")
			break
		}
	}

	// Scrub data
	status.ScrubRunning, status.ScrubErrors, err = status.ScrubData()
	if err != nil {
		return nil, err
	}
	if status.ScrubErrors > 0 {
		failures = append(failures, "scrub errors detected")
	}

	for _, healthRAW := range status.HealthRAW {
		if healthRAW.Temperature == -99 {
			failures = append(failures, "temperature of /dev/"+healthRAW.Device+" cannot be detected")
		}
		// What I have able to learn it seems that temperature and power cycles don't have much effect of life span
		// of a drive so we want to know only about cases the temperature is clearly out of the range of normal
		// operational range.
		if healthRAW.Temperature >= 50 {
			failures = append(failures, "temperature of /dev/"+healthRAW.Device+" is higher 50 °C ("+strconv.Itoa(healthRAW.Temperature)+" °C)")
		}
	}

	// If empty, we will put OK into health field
	status.Health = "OK"
	if len(failures) > 0 {
		status.Health = "PROBLEM"
	}
	status.Failures = failures

	// Utilization stats
	status.UsedBytes, status.TotalBytes, err = status.Df()
	if err != nil {
		return nil, err
	}

	return &status, nil
}
