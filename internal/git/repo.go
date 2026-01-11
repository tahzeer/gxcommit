package git

import "os/exec"

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	return cmd.Run() == nil
}
