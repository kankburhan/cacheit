//go:build darwin

package detector

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
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
	var stat syscall.Stat_t
	if err := syscall.Fstat(int(os.Stdin.Fd()), &stat); err != nil {
		return 0, fmt.Errorf("fstat error: %w", err)
	}
	return uint64(stat.Ino), nil
}

func findCommand(inode uint64) (string, error) {
	// Use different lsof flags for macOS
	cmd := exec.Command("lsof",
		"-Fc",                                    // Output command name
		"-d0",                                    // STDIN file descriptor
		"-a",                                     // AND conditions
		fmt.Sprintf("-iTCP@localhost:%d", inode), // Match inode
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		// Handle "no matches" case gracefully
		if strings.Contains(string(out), "no file use located") {
			return "", fmt.Errorf("no process found writing to pipe")
		}
		return "", fmt.Errorf("lsof error: %v\nOutput: %s", err, string(out))
	}

	return parseLSOF(string(out))
}

func parseLSOF(output string) (string, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "c") {
			cmd := strings.TrimPrefix(line, "c")
			if cmd != "" {
				return cmd, nil
			}
		}
	}
	return "", fmt.Errorf("could not parse command from lsof output:\n%s", output)
}
