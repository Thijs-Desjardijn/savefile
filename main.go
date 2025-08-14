package savefile

import (
	"encoding/gob"
	"encoding/json"
	"errors"

	//"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type EncoderDecoder interface {
	Encode(w io.Writer, v any) error
	Decode(r io.Reader, v any) error
}

type GobCodec struct{}

func (g GobCodec) Encode(w io.Writer, v any) error { return gob.NewEncoder(w).Encode(v) }
func (g GobCodec) Decode(r io.Reader, v any) error { return gob.NewDecoder(r).Decode(v) }

type JSONCodec struct{}

func (j JSONCodec) Encode(w io.Writer, v any) error { return json.NewEncoder(w).Encode(v) }
func (j JSONCodec) Decode(r io.Reader, v any) error { return json.NewDecoder(r).Decode(v) }

type Saver struct {
	dir            string
	maxStoredFiles int
	codec          EncoderDecoder
}

func (s *Saver) fileExt() string {
	switch s.codec.(type) {
	case JSONCodec:
		return ".json"
	case GobCodec:
		return ".bin"
	default:
		return ".dat"
	}
}

// New creates (or reuses) a save directory at path.
func New(path string, codec EncoderDecoder) (*Saver, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, err
	}
	return &Saver{dir: absPath, codec: codec}, nil
}

func NewLimit(path string, codec EncoderDecoder, maxFiles int) (*Saver, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, err
	}
	return &Saver{dir: absPath, maxStoredFiles: maxFiles, codec: codec}, nil
}

// Save writes data to a new file with a timestamp in its name.
func (s *Saver) Save(data any) error {
	s.DeleteOld()
	filename := "save_" + time.Now().Format("20060102_150405") + s.fileExt()
	path := filepath.Join(s.dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return s.codec.Encode(file, data)
}

func (s *Saver) DeleleFile(fileName string) error {
	fullpath := filepath.Join(s.dir, fileName)
	_, err := os.Stat(fullpath)
	if err != nil {
		return err
	}
	err = os.Remove(fullpath)
	if err != nil {
		return err
	}
	return nil
}

func (s *Saver) DeleteOld() error {
	if s.maxStoredFiles != 0 {
		files, err := os.ReadDir(s.dir)
		if err != nil {
			return err
		}
		oldestTime := time.Now()
		var oldestFile string
		saveFilesCount := 0
		for {
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
				saveFilesCount++
				if saveFilesCount == 1 || t.Before(oldestTime) {
					oldestTime = t
					oldestFile = f.Name()
				}
			}
			if saveFilesCount >= s.maxStoredFiles {
				if oldestFile != "" {
					err = os.Remove(filepath.Join(s.dir, oldestFile))
					if err != nil {
						return err
					}
				}
			} else {
				break
			}
		}
	}
	return nil
}

func (s *Saver) Load(file string, target any) error {
	f, err := os.Open(filepath.Join(s.dir, file))
	if err != nil {
		return err
	}
	return s.codec.Decode(f, target)
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
	return s.codec.Decode(file, target)
}
