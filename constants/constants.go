package constants

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mayron1806/gengo/config"
)

type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	Mac     OS = "darwin"
)

func CheckConstants() error {
	_, err := getOS()
	if err != nil {
		return err
	}
	_, err = getGengoHome()
	if err != nil {
		return err
	}
	return nil
}
func getGengoHome() (string, error) {
	gengoHome := os.Getenv("GENGO_HOME")

	if gengoHome == "" {
		config.GetLogger().Warnf("GENGO_HOME not set, trying to use user home directory")
		home, err := os.UserHomeDir()
		if err != nil {
			config.GetLogger().Errorf("Error getting user home directory: %s", home)
			return "", err
		}
		gengoHome = filepath.Join(home, ".gengo")
		err = os.MkdirAll(gengoHome, os.ModePerm)
		if err != nil {
			config.GetLogger().Errorf("Error creating gengo home directory: %s", gengoHome)
			return "", err
		}
		config.GetLogger().Warnf("Using %s, set GENG0_HOME to change", gengoHome)
	}
	return gengoHome, nil
}
func getOS() (OS, error) {
	switch {
	case runtime.GOOS == "windows":
		return Windows, nil
	case runtime.GOOS == "linux":
		return Linux, nil
	case runtime.GOOS == "darwin":
		return Mac, nil
	default:
		return Linux, errors.New("Unsupported OS")
	}
}

func GetGengoHome() (string, error) {
	return getGengoHome()
}
func GetOS() OS {
	os, err := getOS()
	if err != nil {
		config.GetLogger().Warn(err)
	}
	return os
}
