package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func UnmarshalYAMLConfigFile(filePathEnv string, out interface{}) error {
	filePath := os.Getenv(filePathEnv)
	configData, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return errors.Wrapf(err, "read file with file path = %s", filePath)
	}

	if err = yaml.UnmarshalStrict(configData, out); err != nil {
		return errors.Wrap(err, "unmarshal input config data strict")
	}

	return nil
}
