package main

import (
	"flag"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"
)

var FILE_CACHE = map[string]time.Time{}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	command := flag.String("c", "", "Command to execute if a change to the watched files was detected")
	files := flag.String("f", "", "File extensions to watch for seperated by ;, leave empty for all")
	timeout := flag.Duration("d", time.Millisecond*100, "Timeout for checking all files for changes")
	flag.Parse()
	if *command == "" {
		slog.Error("Command Missing")
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
			slog.Info("Change detected, triggering command", "file", file)
			err := execute(*command)
			if err != nil {
				slog.Error("Failed to execute command", "err", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		c := 0
		defer wg.Done()
		for {
			for _, file := range walkDir(".", extensions) {
				if f, ok := FILE_CACHE[file.Path]; !ok {
					FILE_CACHE[file.Path] = file.LastMod
					c++
				} else if file.LastMod.After(f) {
					hasChange <- file.Path
					FILE_CACHE[file.Path] = file.LastMod
				}
			}
			if c > 0 {
				slog.Info("Watching for changes", "files", c, "command", *command, "timeout", *timeout)
				c = 0
			}
			time.Sleep(*timeout)
		}
	}()

	wg.Wait()
}
