package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

const TIMEOUT = time.Millisecond * 50

var FILE_CACHE = map[string]time.Time{}

func main() {
	command := flag.String("c", "", "Command to execute if a change to the watched files was detected")
	files := flag.String("f", "", "File extensions to watch for seperated by ;, leave empty for all")
	flag.Parse()
	if *command == "" {
		slog.Error("Invalid command", "err", "Missing command, please specifiy")
		os.Exit(1)
	}
	extensions := make([]string, 0, 2)
	if *files == "" {
		*files = "."
	} else {
		extensions = strings.Split(*files, ";")
	}

	hasChange := make(chan string, 1)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			file := <-hasChange
			slog.Info("Detected change, triggering command", "filename", file, "command", *command)
			err := execute(*command)
			if err != nil {
				slog.Error("Failed to execute command", "err", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			for _, file := range walkDir(".", extensions) {
				if f, ok := FILE_CACHE[file.Path]; !ok {
					FILE_CACHE[file.Path] = file.LastMod
				} else if file.LastMod.After(f) {
					hasChange <- file.Path
					FILE_CACHE[file.Path] = file.LastMod
				}
			}
			time.Sleep(TIMEOUT)
		}
	}()

	wg.Wait()
}
