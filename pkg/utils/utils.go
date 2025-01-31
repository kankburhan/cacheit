package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func SafePath(root, path string) (string, error) {
	cleanPath := filepath.Clean(path)
	absPath := filepath.Join(root, cleanPath)

	// Verify the resulting path is within the root directory
	if !strings.HasPrefix(absPath, filepath.Clean(root)+string(os.PathSeparator)) {
		return "", os.ErrInvalid
	}

	return absPath, nil
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ValidateUUID(input string) bool {
	_, err := uuid.Parse(input)
	return err == nil
}

func ReadAllSecure(r io.Reader, max int64) ([]byte, error) {
	limitedReader := io.LimitReader(r, max)
	return io.ReadAll(limitedReader)
}
