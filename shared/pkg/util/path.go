package util

import (
	"errors"
	"path/filepath"
	"runtime"
)

// The __filename equivalent
func Filename() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get the current filename")
	}

	return filename, nil
}

// The __dirname equivalent
func Dirname() (string, error) {
	filename, err := Filename()
	if err != nil {
		return "", err
	}

	return filepath.Dir(filename), nil
}
