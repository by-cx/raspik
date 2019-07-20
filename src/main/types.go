package main

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
		Name      string `yaml:"name" json:"name"`
		UUID      string `yaml:"uuid" json:"uuid"`
		Encrypted bool   `yaml:"encrypted" json:"encrypted"`
	} `yaml:"drives" json:"drives"`

	Shares []struct {
		Name  string `yaml:"name" json:"name"`
		Drive uint   `yaml:"drive" json:"drive"`
	} `yaml:"shares" json:"shares"`
}
