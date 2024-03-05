package main

import (
	"flag"

	"log/slog"
)

func main() {
	command := flag.String("c", "", "Command to execute if a change to the watched files was detected")
	files := flag.String("f", "", "Files to watch by extension, if empty watches directory")
	flag.Parse()
	slog.Info("Cli arguments", "command", *command, "files", *files)
}
