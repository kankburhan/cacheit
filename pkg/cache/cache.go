package cache

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kankburhan/cacheit/pkg/utils"
)

var (
	ErrInvalidID   = errors.New("invalid cache ID")
	ErrCacheMiss   = errors.New("cache miss")
	ErrInvalidPath = errors.New("invalid path")
)

type Manager struct {
	CacheDir string
}

type Entry struct {
	ID       string    `json:"id"`
	Label    string    `json:"label"`
	Created  time.Time `json:"created"`
	LastUsed time.Time `json:"last_used"`
	Size     int       `json:"size"`
}

func NewManager() (*Manager, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	dir := filepath.Join(cacheDir, "cacheit")

	if err := os.MkdirAll(filepath.Join(dir, "data"), 0755); err != nil {
		return nil, err
	}

	return &Manager{CacheDir: dir}, nil
}

func (m *Manager) Save(label string, data []byte) (string, error) {
	id := uuid.New().String()

	dataPath, err := utils.SafePath(m.CacheDir, filepath.Join("data", id+".data"))
	if err != nil {
		return "", ErrInvalidPath
	}

	if err := os.WriteFile(dataPath, data, 0644); err != nil {
		return "", err
	}

	entry := Entry{
		ID:       id,
		Label:    label,
		Created:  time.Now().UTC(),
		LastUsed: time.Now().UTC(),
		Size:     len(data),
	}

	if err := m.updateMetadata(entry); err != nil {
		os.Remove(dataPath)
		return "", err
	}

	return id, nil
}

func (m *Manager) Retrieve(id string) ([]byte, error) {
	if !utils.ValidateUUID(id) {
		return nil, ErrInvalidID
	}

	dataPath, err := utils.SafePath(m.CacheDir, filepath.Join("data", id+".data"))
	if err != nil {
		return nil, ErrInvalidPath
	}

	data, err := os.ReadFile(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrCacheMiss
		}
		return nil, err
	}

	if err := m.updateLastUsed(id); err != nil {
		return data, nil // Return data even if metadata update fails
	}

	return data, nil
}

func (m *Manager) ClearAll() error {
	if err := os.RemoveAll(m.CacheDir); err != nil {
		return err
	}
	return os.MkdirAll(filepath.Join(m.CacheDir, "data"), 0755)
}

func (m *Manager) ClearOne(id string) error {
	if !utils.ValidateUUID(id) {
		return ErrInvalidID
	}

	dataPath, err := utils.SafePath(m.CacheDir, filepath.Join("data", id+".data"))
	if err != nil {
		return ErrInvalidPath
	}

	if err := os.Remove(dataPath); err != nil {
		return err
	}

	return m.removeMetadataEntry(id)
}
