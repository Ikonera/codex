package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ikonera/codex/internal/utils"
	"github.com/ikonera/codex/pkg/models"
	"gopkg.in/yaml.v3"
)

const (
	codexDir  = "codex"
	codexFile = "config.yaml"
)

func CheckForConfigFile() error {
	if _, err := os.ReadFile(getCodexConfigPath()); err != nil {
		r := bufio.NewReader(os.Stdin)
		prompter := utils.NewCLIPrompter(r)
		for {
			switch prompter.Prompt("Codex configuration not found, would you initialize it ? [y/n] ") {
			case "y":
				if err := initializeConfigFile(); err != nil {
					log.Fatalf("Can't initialize codex configuration file: %s\n", err.Error())
				}
			case "n":
				fmt.Printf("To configure Codex, create or edit: %s\n", getCodexConfigPath())
				return nil
			default:
				fmt.Print("Please enter valid input ! ")
			}
		}
	}
	return nil
}

func getCodexConfigPath() string {
	return filepath.Join(getUserConfigDir(), codexDir, codexFile)
}

func getUserConfigDir() string {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("User config directory not found: %s\n", err.Error())
	}
	return userConfDir
}

func initializeConfigFile() error {
	if err := os.Mkdir(filepath.Join(getUserConfigDir(), codexDir), 0755); err != nil {
		return err
	}
	configFile, err := os.Create(getCodexConfigPath())
	if err != nil {
		return err
	}
	defer configFile.Close()

	defaultConf, err := setDefaultConfig()
	if err != nil {
		return nil
	}
	if err := writeConfig(defaultConf, configFile); err != nil {
		return err
	}
	userConf, err := askForUserConfig()
	if err != nil {
		return err
	}
	if err := writeConfig(userConf, configFile); err != nil {
		return err
	}

	return nil
}

func askForUserConfig() (*models.Config, error) {
	r := bufio.NewReader(os.Stdin)
	prompter := utils.NewCLIPrompter(r)
	userAk := prompter.Prompt("Enter a valid access key: ")
	userSk := prompter.Prompt("Enter a valid secret key: ")
	userCredentials := models.NewCredentials(userAk, userSk)
	fmt.Println("Register at least one Codex.")
	codexName := prompter.Prompt("Codex name: ")
	codexSource := prompter.Prompt("Codex source: ")
	codexBucket := prompter.Prompt("Codex destination bucket: ")
	userCodexes := []*models.Codex{models.NewCodex(codexName, codexSource, codexBucket)}
	userConfig := models.NewConfig(userCredentials, userCodexes)
	return userConfig, nil
}

func setDefaultConfig() (*models.Config, error) {
	defaultCredentials := models.NewCredentials("", "")
	defaultCodexes := []*models.Codex{models.NewCodex("default", "", "")}
	defaultConfig := models.NewConfig(defaultCredentials, defaultCodexes)

	return defaultConfig, nil
}

func writeConfig(config *models.Config, file *os.File) error {
	yamlConfig, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if _, err := file.Write(yamlConfig); err != nil {
		return err
	}
	fmt.Println("Saved config !")
	return nil
}
