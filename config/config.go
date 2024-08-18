package config

type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	Mac     OS = "darwin"
)

var logger *Logger

func GetLogger() *Logger {
	if logger != nil {
		return logger
	}

	logger = NewLogger()
	return logger
}
