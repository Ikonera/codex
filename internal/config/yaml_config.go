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

type YAMLConfigManager struct{}

func NewYAMLConfigManager() *YAMLConfigManager {
	return &YAMLConfigManager{}
}

func (m *YAMLConfigManager) CheckForConfig() error {
	if _, err := os.ReadFile(getConfigPath()); err != nil {
		r := bufio.NewReader(os.Stdin)
		prompter := utils.NewCLIPrompter(r)
		for {
			switch prompter.Prompt("Codex configuration not found, would you initialize it ? [y/n] ") {
			case "y":
				if err := m.InitializeConfig(); err != nil {
					log.Fatalf("Can't initialize codex configuration file: %s\n", err.Error())
				}
				return nil
			case "n":
				fmt.Printf("To configure Codex, create or edit: %s\n", getConfigPath())
				return nil
			default:
				fmt.Print("Please enter valid input ! ")
			}
		}
	}
	return nil
}

func getUserConfigDir() string {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("User config directory not found: %s\n", err.Error())
	}
	return userConfDir
}

func getConfigPath() string {
	return filepath.Join(getUserConfigDir(), codexDir, codexFile)
}

func (m *YAMLConfigManager) InitializeConfig() error {
	if err := os.Mkdir(filepath.Join(getUserConfigDir(), codexDir), 0755); err != nil {
		return err
	}
	configFile, err := os.Create(getConfigPath())
	if err != nil {
		return err
	}
	defer configFile.Close()

	var config *models.Config

	config = setDefaultConfig()

	userConf, err := askForUserConfig()
	if err != nil {
		return err
	}
	config = userConf
	if err := m.WriteConfig(config); err != nil {
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
	fmt.Println("\nRegister at least one Codex.")
	codexName := prompter.Prompt("Codex name: ")

	codexSource := registerLocalSource(prompter)

	codexBucket := prompter.Prompt("Codex destination bucket: ")
	userCodexes := []*models.Codex{models.NewCodex(codexName, codexSource, codexBucket)}
	userConfig := models.NewConfig(userCredentials, userCodexes)
	return userConfig, nil
}

func registerLocalSource(p *utils.CLIPrompter) string {
	var path string
	for {
		path = p.Prompt("Codex local source: ")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Print("Please enter a valid local path ! ")
		} else {
			break
		}
	}
	return path
}

func setDefaultConfig() *models.Config {
	defaultCredentials := models.NewCredentials("", "")
	defaultCodexes := []*models.Codex{models.NewCodex("default", "", "")}
	defaultConfig := models.NewConfig(defaultCredentials, defaultCodexes)

	return defaultConfig
}

func (m *YAMLConfigManager) ReadConfig() (*models.Config, error) {
	fileContent, err := os.ReadFile(getConfigPath())
	if err != nil {
		return nil, err
	}
	var config *models.Config
	if err := yaml.Unmarshal(fileContent, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func (m *YAMLConfigManager) WriteConfig(config *models.Config) error {
	yamlConfig, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	file, err := os.Open(getConfigPath())
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(yamlConfig); err != nil {
		return err
	}
	fmt.Println("Saved config !")
	return nil
}
