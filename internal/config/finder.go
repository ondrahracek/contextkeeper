package config

import (
	"os"
	"path/filepath"
)

// Finder finds configuration files
type Finder struct{}

func NewFinder() *Finder {
	return &Finder{}
}

func (f *Finder) FindConfig() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, ".contextkeeper"), nil
}

func (f *Finder) FindUpwards(filename string) (string, error) {
	return "", nil
}
