package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Ext int

const (
	YAML Ext = iota
	JSON
	ENV
)

const (
	githubAppIdentifier = "GITHUB_APP_IDENTIFIER"
	githubClientSecret  = "GITHUB_CLIENT_SECRET"
	githubPrivateKey    = "GITHUB_PRIVATE_KEY"
	githubWebhookSecret = "GITHUB_WEBHOOK_SECRET"
)

type Config struct {
	GHAppId         int
	GHClientSecret  string
	GHPrivateKey    *rsa.PrivateKey
	GHWebhookSecret string
}

func (e Ext) String() string {
	switch e {
	case YAML:
		return "yaml"
	case JSON:
		return "json"
	case ENV:
		return "env"
	default:
		return ""
	}
}

func Load(filename string, paths []string, extension Ext) (Config, error) {
	viper.SetConfigName(filename)
	viper.SetConfigType(extension.String())
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("could not read config file {file_name: %s.%s, error: %v}", filename, extension, err)
	}

	ghAppId, err := loadInt(githubAppIdentifier)
	if err != nil {
		return Config{}, err
	}

	ghClientSecret, err := load(githubClientSecret)
	if err != nil {
		return Config{}, err
	}

	ghPrivateKey, err := loadPrivateKey(githubPrivateKey)
	if err != nil {
		return Config{}, err
	}

	ghWebhookSecret, err := load(githubWebhookSecret)
	if err != nil {
		return Config{}, err
	}

	return Config{GHAppId: ghAppId, GHClientSecret: ghClientSecret, GHPrivateKey: ghPrivateKey, GHWebhookSecret: ghWebhookSecret}, nil
}

func loadInt(key string) (int, error) {
	if !viper.InConfig(githubAppIdentifier) {
		return -1, fmt.Errorf("could not get %s", key)
	}
	return viper.GetInt(key), nil
}

func loadPrivateKey(key string) (*rsa.PrivateKey, error) {
	rawPrivateKey, err := load(key)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(rawPrivateKey))
	if block == nil {
		return nil, errors.New("could not decode RSA Private Key")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("RSA Private Key is the wrong type. Expected RSA PRIVATE KEY, but got %s", block.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("could not parse RSA Private Key: %w", err)
	}

	return privateKey, nil
}

func load(key string) (string, error) {
	if !viper.InConfig(key) {
		return "", fmt.Errorf("could not get %s", key)
	}
	return viper.GetString(key), nil
}

func (c Config) String() string {
	return fmt.Sprintf("<Config {%s: %d, %s: %s, %s: %d, %s: %s}>", githubAppIdentifier, c.GHAppId, githubClientSecret, c.GHClientSecret, githubPrivateKey, c.GHPrivateKey.N, githubWebhookSecret, c.GHWebhookSecret)
}
