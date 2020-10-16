package main

const jsonIdent = "  "

// Config represents configuration of the NAS environment
type Config struct {
	General struct {
		SambaHomeDirectories bool   `yaml:"samba_home_directories" json:"samba_home_directories"`
		HomesDrive           uint   `yaml:"homes_drive" json:"homes_drive"`
		SharedGroup          string `yaml:"shared_group" json:"shared_group"`
	} `yaml:"general" json:"general"`

	Backup struct {
		KeepDaily   uint              `yaml:"keep_daily" json:"keep_daily"`
		KeepWeekly  uint              `yaml:"keep_weekly" json:"keep_weekly"`
		KeepMonthly uint              `yaml:"keep_monthly" json:"keep_monthly"`
		ResticEnv   map[string]string `yaml:"restic_env" json:"restic_env"`
	} `yaml:"backup" json:"backup"`

	Users []struct {
		Name         string `yaml:"name" json:"name"`
		Password     string `yaml:"password" json:"-"`
		PasswordHash string `yaml:"password_hash" json:"password_hash"`
		Services     struct {
			Syncthing struct {
				Enabled bool `yaml:"enabled" json:"enabled"`
			} `yaml:"syncthing" json:"syncthing"`
		} `yaml:"services" json:"services"`
	}

	Drives []struct {
		Name        string   `yaml:"name" json:"name"`
		UUID        string   `yaml:"uuid" json:"uuid"`
		Encrypted   bool     `yaml:"encrypted" json:"encrypted"`
		Raid        bool     `yaml:"raid" json:"raid"` // Only Btrfs raid is supported
		RaidDevices []string `yaml:"raid_devices" json:"raid_devices"`
	} `yaml:"drives" json:"drives"`

	Shares []struct {
		Name  string `yaml:"name" json:"name"`
		Drive uint   `yaml:"drive" json:"drive"`
	} `yaml:"shares" json:"shares"`
}

// DeviceStatus contains information about status of the local system
type DeviceStatus struct {
	CPUTemperature float64 `json:"cpu_temperature"`
}

// DrivesStatus contains data about health of all drives in the system
type DrivesStatus struct {
	SMARTHealth      Health           `json:"smart_health"`
	FileSystemErrors FileSystemErrors `json:"filesystem_errors"`
}

// Health contains a few metrics about current health of one physical drive
type Health struct {
	Device           string `json:"device"`
	SMARTHealth      bool   `json:"smart_health"`
	RelocatedSectors int    `json:"relocated_sectors"`
	PendingSectors   int    `json:"pending_sectors"`
	Temperature      int    `json:"temperature"`
}

// FileSystemErrors contains number of errors btrfs detected
type FileSystemErrors struct {
	WriteIOErrors    int `json:"write_io_errs"`
	ReadIOErrors     int `json:"read_io_errs"`
	FlushIOErrors    int `json:"flush_io_errs"`
	CorruptionErrors int `json:"corruption_errs"`
	GenerationErrsor int `json:"generation_errs"`
}
