package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	SessionID      string    `json:"session_id"`
	ExpirationTime time.Time `json:"expiration_time"`
}

func configPath() (string, error) {

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	// store the current session of user in ~/.login-cli/config.json

	return filepath.Join(
		homeDir,
		".login-cli",
		"config.json",
	), nil
}

func Load() (*Config, error) {

	path, err := configPath()
	if err != nil {
		return nil, err
	}

	// read the config file of current session
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// deserialize the session content into config object
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// this will overrite the config (session) file if it already exists else create a new session file
func Save(sessionId string, expirationTime time.Time) error {

	path, err := configPath()

	if err != nil {
		return err
	}

	// create a config file if it does not exist
	err = os.MkdirAll(filepath.Dir(path), 0755)

	if err != nil {
		return err
	}

	// create a config object
	config := Config{
		SessionID:      sessionId,
		ExpirationTime: expirationTime,
	}
	data, err := json.MarshalIndent(config, "", " ")

	if err != nil {
		return err
	}

	// overwrite the file with current config data
	return os.WriteFile(path, data, 0644)
}

func Delete() error {

	path, err := configPath()

	if err != nil {
		return err
	}

	// check if path exists , if it does not exist (not do anything)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	// if it does remove the config file
	return os.WriteFile(path, []byte("{}"), 0644)
}
