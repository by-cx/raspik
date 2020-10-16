package main

import (
	"log"
	"os/exec"
	"strings"
)

func runCommand(cmd string, args []string, stdin []byte) ([]byte, error) {
	log.Println("Command:", cmd, strings.Join(args, " "))

	subprocess := exec.Command(cmd, args...)

	if stdin != nil {
		writer, err := subprocess.StdinPipe()
		if err != nil {
			return nil, err
		}
		_, err = writer.Write(stdin)
		if err != nil {
			return nil, err
		}
	}

	out, err := subprocess.Output()
	return out, err
}

func startService(name string) error {
	_, err := runCommand("/bin/systemctl", []string{"start", name}, nil)
	return err
}
