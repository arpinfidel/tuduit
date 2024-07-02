package config

import (
	"encoding/json"
	"os"
)

func NewFromPath(path string) (cfg Configuration, err error) {
	jsonFile, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(jsonFile, &cfg.Config)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

type SecretManager interface {
	GetSecret() (secret []byte, err error)
}

func (c *Configuration) LoadSecret(sm SecretManager) (err error) {
	s, err := sm.GetSecret()
	if err != nil {
		return err
	}

	defer c.Lock()()

	return json.Unmarshal(s, &c.Config)
}
