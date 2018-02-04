package server

import (
	"io"
	"os"
	"path/filepath"
)

// Server represents the entity that receives files
type Server struct {
	Root    string
	Fridge  string
	Freezer string
}

// New creates a Server
func New(root string) Server {
	fridge := filepath.Join(root, "fridge")
	freezer := filepath.Join(root, "freezer")
	return Server{root, fridge, freezer}
}

// Remove file from server
func (s Server) Remove(path string) error {
	dest := filepath.Join(s.Fridge, filepath.Base(path))
	return os.RemoveAll(dest)
}

// Add file to server
func (s Server) Add(path string) error {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	outpath := filepath.Join(s.Fridge, filepath.Base(path))
	fileinfo, err := in.Stat()
	if err != nil {
		return err
	}
	if fileinfo.IsDir() {
		return os.Mkdir(outpath, os.ModePerm)
	}
	return copyFile(outpath, in)
}

func copyFile(path string, in *os.File) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}
