package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

const purgeInterval = (15 * time.Minute)

func remove(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Printf("Error when removing %s: %v\n", path, err)
	}
}

func isHiddenFile(name string) bool {
	return len(name) > 0 && name[0] == '.'
}
func purgeOldEntries(directory string, limit time.Time) {
	infos, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Printf("Error when purging directory %s: %v\n", directory, err)
		return
	}

	for _, entry := range infos {
		if !entry.IsDir() && !isHiddenFile(entry.Name()) && entry.ModTime().Before(limit) {
			path := path.Join(directory, entry.Name())
			remove(path)
		}
	}
}

func purgeWorker(directory string, purgeOlder time.Duration) {
	for {
		limit := time.Now().Add(-purgeOlder)
		purgeOldEntries(directory, limit)
		time.Sleep(purgeInterval)

	}
}

func initPurge(directory string, purgeOlder time.Duration) {
	go purgeWorker(directory, purgeOlder)
}
