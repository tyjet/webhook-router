package config

import (
	"github.com/spf13/viper"
	log "github.com/tyjet/soos-ult"
	"go.uber.org/zap"
)

type kafka struct {
	bootstrap *bootstrap
}

type bootstrap struct {
	servers string
}

type config struct {
	kafka *kafka
}

var cfg *config

type Type int

const (
	DOTENV Type = iota
	ENV
	JSON
	PROP
	PROPERTIES
	PROPS
	TOML
	YAML
	YML
)

func (t Type) String() string {
	switch t {
	case DOTENV:
		return "dotenv"
	case ENV:
		return "env"
	case JSON:
		return "json"
	case PROP:
		return "prop"
	case PROPERTIES:
		return "properties"
	case PROPS:
		return "props"
	case TOML:
		return "toml"
	case YAML:
		return "yaml"
	case YML:
		return "yml"
	default:
		log.Error("unknown type", zap.Int("config_type", int(t)))
		return "unknown"
	}
}

func NewConfig(name string, cfgType Type, paths []string) (*config, error) {
	viper.SetConfigName(name)
	viper.SetConfigType(cfgType.String())
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Error("could not read config file", zap.Error(err), zap.String("config_name", name), zap.String("config_type", cfgType.String()))
		return nil, err
	}

	if !viper.InConfig("kafka.bootstrap.servers") {
		log.Error("could not read kafka.bootstrap.servers")
		return nil, err
	}
	bootstrapServers := viper.GetString("kafka.bootstrap.servers")

	bootstrap := &bootstrap{servers: bootstrapServers}
	kafka := &kafka{bootstrap: bootstrap}
	return &config{kafka: kafka}, nil
}

func ReadConfig(name string, cfgType Type, paths []string) bool {
	conf, err := NewConfig(name, cfgType, paths)
	if err != nil {
		return false
	}

	cfg = conf
	return true
}

func KafkaBootstrapServers() string {
	return cfg.kafka.bootstrap.servers
}
