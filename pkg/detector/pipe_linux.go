//go:build linux

package detector

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func (d *PipeDetector) DetectCommand() (string, error) {
	if !isPiped() {
		return "", ErrNotPipe
	}

	inode, err := getInode()
	if err != nil {
		return "", err
	}
	return findCommand(inode)
}

func getInode() (uint64, error) {
	target, err := os.Readlink("/proc/self/fd/0")
	if err != nil {
		return 0, fmt.Errorf("readlink error: %w", err)
	}

	if !strings.HasPrefix(target, "pipe:[") {
		return 0, ErrNotPipe
	}

	inodeStr := strings.TrimPrefix(target, "pipe:[")
	inodeStr = strings.TrimSuffix(inodeStr, "]")
	return strconv.ParseUint(inodeStr, 10, 64)
}

func findCommand(inode uint64) (string, error) {
	cmd := exec.Command("lsof", "-t", "-c", "bash", "-a", "-d0", "-Fn")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("lsof error: %w", err)
	}

	return parseLSOF(string(out), inode)
}

func parseLSOF(output string, targetInode uint64) (string, error) {
	lines := strings.Split(output, "\n")
	var currentPid string

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "p"):
			currentPid = strings.TrimPrefix(line, "p")
		case strings.HasPrefix(line, "n"):
			fd := strings.TrimPrefix(line, "n")
			if strings.Contains(fd, "pipe:") {
				inodeStr := strings.TrimPrefix(fd, "pipe:[")
				inodeStr = strings.TrimSuffix(inodeStr, "]")
				inode, err := strconv.ParseUint(inodeStr, 10, 64)
				if err != nil {
					continue
				}
				if inode == targetInode && currentPid != "" {
					return getCommandByPID(currentPid)
				}
			}
		}
	}
	return "", fmt.Errorf("command not found for inode %d", targetInode)
}

func getCommandByPID(pid string) (string, error) {
	cmd := exec.Command("ps", "-o", "command=", "-p", pid)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ps error: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
