package main

import (
	"bytes"
	"strconv"
)

// RPi encapsulates methods to gather info from VC commands
type RPi struct{}

// CPUTemperature returns temperature of CPU in °C
func (r *RPi) CPUTemperature() (float64, error) {
	// Output example: temp=45.0'C
	out, err := runCommand("vcgencmd", []string{"measure_temp"}, []byte(""))
	if err != nil {
		return 0, err
	}

	out = bytes.Replace(out, []byte("temp="), []byte(""), 1)
	out = bytes.Replace(out, []byte("'C\n"), []byte(""), 1)
	value, err := strconv.ParseFloat(string(out), 64)
	return value, err
}
