package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

type Config struct {
	APIVersion string `yaml:"apiVersion" validate:"required"`
	Spec       Spec   `yaml:"spec" validate:"required"`
}

type Spec struct {
	System    System        `yaml:"system" validate:"required"`
	Storages  []SpecStorage `yaml:"storages" validate:"required"`
	Databases []Databases   `yaml:"databases" validate:"required"`
	Backups   []Backup      `yaml:"backups" validate:"required"`
}

type Backup struct {
	Name      string             `yaml:"name" validate:"required"`
	Schedule  string             `yaml:"schedule" validate:"required"`
	Databases []DatabasesElement `yaml:"databases" validate:"required"`
	Storage   SpecStorageName    `yaml:"storage" validate:"required"`
}

type DatabasesElement struct {
	Name string `yaml:"name" validate:"required"`
}

type Databases struct {
	Name     string `yaml:"name" validate:"required"`
	Type     string `yaml:"type" validate:"required"`
	Host     string `yaml:"host" validate:"required"`
	Port     int64  `yaml:"port" validate:"required"`
	User     string `yaml:"user" validate:"required"`
	Password string `yaml:"password" validate:"required"`
}

type SpecStorage struct {
	SpecStorageName
	SpecStorageS3
}

type SpecStorageName struct {
	Name string `yaml:"name" validate:"required"`
}

type SpecStorageS3 struct {
	S3 S3 `yaml:"s3" validate:"required"`
}

type S3 struct {
	Region       string `yaml:"region" validate:"required"`
	Bucket       string `yaml:"bucket" validate:"required"`
	AccessKey    string `yaml:"access-key" validate:"required"`
	ClientSecret string `yaml:"client-secret" validate:"required"`
}

type System struct {
	LogLevel string `yaml:"logLevel" validate:"required"`
	Web      Web    `yaml:"web" validate:"required"`
}

type Web struct {
	Port int `yaml:"port" validate:"required"`
}

func NewConfig(configPath string) (*Config, error) {
	validate := validator.New()
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	err = validate.Struct(config)
	if err != nil {
		return nil, err
	}
	return config, nil

}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "config.yaml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}