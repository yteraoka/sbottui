package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const cacheDir = "sbottui"

// dir returns the cache directory path (~/.cache/sbottui/).
func dir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}
	return filepath.Join(base, cacheDir), nil
}

// Load reads a cached JSON file into v. Returns os.ErrNotExist if not cached.
func Load(name string, v interface{}) error {
	d, err := dir()
	if err != nil {
		return err
	}
	path := filepath.Join(d, name+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err // caller checks os.IsNotExist
	}
	return json.Unmarshal(data, v)
}

// Save atomically writes v as JSON to the cache file name.
func Save(name string, v interface{}) error {
	d, err := dir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(d, 0o755); err != nil {
		return fmt.Errorf("mkdir cache dir: %w", err)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshal cache: %w", err)
	}

	finalPath := filepath.Join(d, name+".json")
	tmp, err := os.CreateTemp(d, name+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()        //nolint:errcheck
		os.Remove(tmpPath) //nolint:errcheck
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmpPath) //nolint:errcheck
		return fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		os.Remove(tmpPath) //nolint:errcheck
		return fmt.Errorf("rename temp file: %w", err)
	}
	return nil
}

// Clear removes the cache file for name.
func Clear(name string) error {
	d, err := dir()
	if err != nil {
		return err
	}
	path := filepath.Join(d, name+".json")
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// ClearAll removes all cache files for sbottui.
func ClearAll() error {
	for _, name := range []string{"devices", "scenes"} {
		if err := Clear(name); err != nil {
			return err
		}
	}
	return nil
}
