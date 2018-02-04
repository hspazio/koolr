package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var watcher *fsnotify.Watcher
var root string

func watchDir(path string, info os.FileInfo, err error) error {
	log.Println(path)
	return nil
}

func removePath(path string) error {
	dest := filepath.Join(root, filepath.Base(path))
	return os.RemoveAll(dest)
}

func addPath(path string) error {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	outpath := filepath.Join(root, filepath.Base(path))
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

func main() {
	root = "/Users/fabio/koolr-server"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Create:
					log.Println("copying file:", event.Name)
					if err := addPath(event.Name); err != nil {
						log.Fatalf("could not add path: %s", err)
					}
				case fsnotify.Write:
					log.Println("copying new version of file:", event.Name)
					if err := addPath(event.Name); err != nil {
						log.Fatalf("could not add path: %s", err)
					}
				// case fsnotify.Rename:
				// 	log.Println("renaming file:", event.Name)
				case fsnotify.Remove:
					log.Println("removing file:", event.Name)
					if err := removePath(event.Name); err != nil {
						log.Fatalf("could not remove path: %s", err)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("/Users/fabio/koolr-client")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
