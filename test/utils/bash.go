package utils

import (
	"os"
	"os/exec"
)

func BashRun(cmds string) error {
	cmd := exec.Command("bash", "-c", cmds)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
