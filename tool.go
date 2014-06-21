package main

import (
	"code.google.com/p/go.exp/fsnotify"
	"log"
	"os"
	"os/exec"
	"strings"
)

func gotool(done chan struct{}, commands ...string) {
	log.Println("go", commands)
	cmd := exec.Command("go", commands...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	done <- struct{}{}
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}

	build := true
	done := make(chan struct{})

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") && build && !ev.IsAttrib() {
				build = false
				go gotool(done, os.Args[1:]...)
			}
		case err := <-watcher.Error:
			log.Fatal("error", err)
		case <-done:
			build = true
		}
	}
}
