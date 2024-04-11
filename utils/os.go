package utils

import (
	"os/exec"
	"strconv"
)

func IsRootPrivilege() (v bool) {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		return
	}

	// 0 = root, 501 = non-root user
	i, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		return
	}

	if i != 0 {
		return
	}

	v = true
	return
}
