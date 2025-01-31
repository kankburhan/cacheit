package detector

import (
	"errors"
	"os"
)

var (
	ErrNotPipe       = errors.New("not a pipe")
	ErrUnsupportedOS = errors.New("unsupported operating system")
)

type PipeDetector struct{}

func NewPipeDetector() *PipeDetector {
	return &PipeDetector{}
}

// Common utility function
func isPiped() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) == 0
}
