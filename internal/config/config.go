package config

import (
	"github.com/ikonera/codex/pkg/models"
)

type IConfigManager interface {
	CheckForConfig() error
	InitializeConfig() error
	ReadConfig() (*models.Config, error)
	WriteConfig(*models.Config) error
}
