package main

// Config represents configuration of the NAS environment
type Config struct {
	General struct {
		SambaHomeDirectories bool   `yaml:"samba_home_directories"`
		HomesDrive           uint   `yaml:"homes_drive"`
		SharedGroup          string `yaml:"shared_group"`
	} `yaml:"general"`

	Backup struct {
		KeepDaily   uint              `yaml:"keep_daily"`
		KeepWeekly  uint              `yaml:"keep_weekly"`
		KeepMonthly uint              `yaml:"keep_monthly"`
		ResticEnv   map[string]string `yaml:"restic_env"`
	} `yaml:"backup"`

	Users []struct {
		Name         string `yaml:"name"`
		Password     string `yaml:"password"`
		PasswordHash string `yaml:"password_hash"`
		Services     struct {
			Syncthing struct {
				Enabled bool `yaml:"enabled"`
			} `yaml:"services"`
		} `yaml:"syncthing"`
	}

	Drives []struct {
		Name      string `yaml:"name"`
		UUID      string `yaml:"uuid"`
		Encrypted bool   `yaml:"encrypted"`
	} `yaml:"drives"`

	Shares []struct {
		Name  string `yaml:"name"`
		Drive uint   `yaml:"drive"`
	} `yaml:"shares"`
}
