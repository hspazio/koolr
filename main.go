package main

import (
	"log"
	"os"

	"github.com/fsnotify/fsnotify"

	"github.com/hspazio/koolr/server"
)

var watcher *fsnotify.Watcher
var svr server.Server

func watchDir(path string, info os.FileInfo, err error) error {
	log.Println(path)
	return nil
}

func main() {
	svr = server.New("/Users/fabio/koolr-server")
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
					if err := svr.Add(event.Name); err != nil {
						log.Fatalf("could not add path: %s", err)
					}
				case fsnotify.Write:
					log.Println("copying new version of file:", event.Name)
					if err := svr.Add(event.Name); err != nil {
						log.Fatalf("could not add path: %s", err)
					}
				// case fsnotify.Rename:
				// 	log.Println("renaming file:", event.Name)
				case fsnotify.Remove:
					log.Println("removing file:", event.Name)
					if err := svr.Remove(event.Name); err != nil {
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
