package config

import (
	"github.com/SwirlGit/queue-scheduler/pkg/config"
	"github.com/SwirlGit/queue-scheduler/pkg/database/postgres"
	"github.com/pkg/errors"
)

type Config struct {
	QSDB postgres.Config `yaml:"qs-db"`
}

func InitConfig(filePath string) (Config, error) {
	var cfg Config
	if err := config.UnmarshalYAMLConfigFile(filePath, &cfg); err != nil {
		return Config{}, errors.Wrap(err, "unmarshal yaml config file")
	}
	return cfg, nil
}
