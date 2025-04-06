package main

import (
	"log"

	"github.com/ikonera/codex/handlers/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		log.Fatalf("Failed to run codex: %s\n", err.Error())
	}
}
