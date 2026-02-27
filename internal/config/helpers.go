package config

import (
	"fmt"
	"os"
)

func parseEnvVariable(env string) (AppEnv, error) {
	switch env {
	case string(Development):
		return Development, nil
	case string(Release):
		return Release, nil
	default:
		return "", ErrInvalidAppEnv
	}
}

func checkConfigPathExists(path string) error {

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %v", ErrConfigPathNotExists, err)
		}
		return err
	}

	return nil
}
