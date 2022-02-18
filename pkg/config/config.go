package config

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func UnmarshalYAMLConfigFile(filePath string, out interface{}) error {
	configData, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return errors.Wrap(err, "read file")
	}

	if err = yaml.UnmarshalStrict(configData, out); err != nil {
		return errors.Wrap(err, "unmarshal input config data strict")
	}

	return nil
}
