package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/kankburhan/cacheit/pkg/utils"
)

func (m *Manager) updateMetadata(entry Entry) error {
	entries, err := m.loadMetadata()
	if err != nil {
		return err
	}

	entries = append(entries, entry)
	return m.saveMetadata(entries)
}

func (m *Manager) loadMetadata() ([]Entry, error) {
	metaPath := filepath.Join(m.CacheDir, "metadata.json")

	if !utils.FileExists(metaPath) {
		return []Entry{}, nil
	}

	file, err := os.Open(metaPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []Entry
	if err := json.NewDecoder(file).Decode(&entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (m *Manager) saveMetadata(entries []Entry) error {
	metaPath := filepath.Join(m.CacheDir, "metadata.json")
	file, err := os.Create(metaPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(entries)
}

func (m *Manager) updateLastUsed(id string) error {
	entries, err := m.loadMetadata()
	if err != nil {
		return err
	}

	for i, entry := range entries {
		if entry.ID == id {
			entries[i].LastUsed = time.Now().UTC()
			break
		}
	}

	return m.saveMetadata(entries)
}

func (m *Manager) removeMetadataEntry(id string) error {
	entries, err := m.loadMetadata()
	if err != nil {
		return err
	}

	newEntries := make([]Entry, 0, len(entries))
	for _, entry := range entries {
		if entry.ID != id {
			newEntries = append(newEntries, entry)
		}
	}

	return m.saveMetadata(newEntries)
}

func (m *Manager) LoadMetadata() ([]Entry, error) {
	return m.loadMetadata()
}
