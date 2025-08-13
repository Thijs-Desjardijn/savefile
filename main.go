package savefile

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"time"
)

type Saver struct {
	dir            string
	maxStoredFiles int
}

// New creates (or reuses) a save directory at path.
func New(path string) (*Saver, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, err
	}
	return &Saver{dir: absPath}, nil
}

func NewLimit(path string, maxFiles int) (*Saver, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, err
	}
	return &Saver{dir: absPath, maxStoredFiles: maxFiles}, nil
}

// Save writes data to a new file with a timestamp in its name.
func (s *Saver) Save(data any) error {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return err
	}
	if s.maxStoredFiles != 0 {
		latestTime := time.Now()
		var latestFile string
		for _, f := range files {
			if !f.Type().IsRegular() {
				continue
			}
			if len(f.Name()) < 20 { // minimal length check for timestamp pattern
				continue
			}
			timestamp := f.Name()[5:20]
			t, err := time.Parse("20060102_150405", timestamp)
			if err != nil {
				continue // skip files that don't match
			}
			if t.Before(latestTime) {
				latestTime = t
				latestFile = f.Name()
			}
		}
		if len(files) >= s.maxStoredFiles {
			err = os.Remove(latestFile)
			if err != nil {
				return err
			}
		}
	}
	filename := "save_" + time.Now().Format("20060102_150405") + ".bin"
	path := filepath.Join(s.dir, filename)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(data)
}

// LoadLatest reads and decodes the most recent save file.
func (s *Saver) LoadLatest(target any) error {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no save files found")
	}

	var mostRecentFile string
	var mostRecentTime time.Time
	for _, f := range files {
		if !f.Type().IsRegular() {
			continue
		}
		if len(f.Name()) < 20 { // minimal length check for timestamp pattern
			continue
		}
		timestamp := f.Name()[5:20]
		t, err := time.Parse("20060102_150405", timestamp)
		if err != nil {
			continue // skip files that don't match
		}
		if t.After(mostRecentTime) {
			mostRecentTime = t
			mostRecentFile = f.Name()
		}
	}
	if mostRecentFile == "" {
		return errors.New("no valid save files found")
	}

	file, err := os.Open(filepath.Join(s.dir, mostRecentFile))
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewDecoder(file).Decode(target)
}
