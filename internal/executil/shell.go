package executil

import (
	"os"
	"os/exec"
)

func RunScript(script string) error {
	tmp := "gxcommit.sh"

	if err := os.WriteFile(tmp, []byte(script), 0755); err != nil {
		return err
	}

	cmd := exec.Command("bash", tmp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
