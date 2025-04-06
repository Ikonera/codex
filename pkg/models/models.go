package models

type Config struct {
	Credentials *Credentials `yaml:"credentials"`
	Codexes     []*Codex     `yaml:"codexes"`
}

type Credentials struct {
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}

type Codex struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
	Bucket string `yaml:"bucket"`
}

func NewCodex(name, source, bucket string) *Codex {
	return &Codex{
		Name:   name,
		Source: source,
		Bucket: bucket,
	}
}

func NewCredentials(ak, sk string) *Credentials {
	return &Credentials{
		AccessKey: ak,
		SecretKey: sk,
	}
}

func NewConfig(creds *Credentials, codexes []*Codex) *Config {
	return &Config{
		Credentials: creds,
		Codexes:     codexes,
	}
}
