package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

const padLen = 3

var counter = 0
var prefix string

func padCounter(counter int) string {
	return fmt.Sprintf("%0*d", padLen, counter)
}

func renameFile(oldPath string) {
	dir := filepath.Dir(oldPath)
	ext := filepath.Ext(oldPath)
	base := filepath.Base(oldPath)
	base = strings.TrimSuffix(base, ext)

	if strings.HasPrefix(base, prefix) {
		return
	}

	newBase := prefix + padCounter(counter)
	newPath := filepath.Join(dir, newBase+ext)

	err := os.Rename(oldPath, newPath)
	if err != nil {
		log.Printf("Failed to rename file %s: %s\n", oldPath, err)
	} else {
		log.Printf("Renamed %s to %s\n", oldPath, newPath)
		counter++
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: [path] <prefix>")
		os.Exit(1)
	}

	var path string
	if len(os.Args) == 2 {
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		path = currentDir
		prefix = os.Args[1]
	} else {
		path = os.Args[1]
		prefix = os.Args[2]
	}

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
				if event.Op&fsnotify.Create == fsnotify.Create {
					// time.Sleep(1 * time.Second) // Ensure file is fully written
					renameFile(event.Name)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Watching directory: %s\n", path)
	<-done
}
