package cli

import (
	"fmt"
	"os"

	"github.com/ikonera/codex/internal/config"
)

var usage = `
Usage:
  codex [command]

Available commands:
  - list: 'codex list' display the configured codexes available in your configuration file.
  - add : 'codex add' prompts you to configure a new codex.
`

func Run() error {
	if len(os.Args) < 2 {
		printUsage()
		return nil
	}

	cfgManager := config.NewYAMLConfigManager()

	switch os.Args[1] {
	case "init":
		if err := handleInitCmd(cfgManager); err != nil {
			return err
		}
	case "help":
		printUsage()
	default:
		printUsage()
	}

	return nil
}

func handleInitCmd(m config.IConfigManager) error {
	displayHeader()

	if err := m.CheckForConfig(); err != nil {
		return err
	}

	return nil
}

func printUsage() {
	fmt.Println(usage)
}

func displayHeader() {
	fmt.Println(`
=============================================

_________            .___
\_   ___ \  ____   __| _/____ ___  ___
/    \  \/ /  _ \ / __ |/ __ \\  \/  /
\     \___(  <_> ) /_/ \  ___/ >    <
 \______  /\____/\____ |\___  >__/\_ \
        \/            \/    \/      \/

Name: Codex
Author: Ikonera
Github repo: https://github.com/ikonera/codex
=============================================
	`)
}
