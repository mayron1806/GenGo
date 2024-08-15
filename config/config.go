package config

type OS string

const (
	Windows OS = "windows"
	Linux   OS = "linux"
	Mac     OS = "darwin"
)

var (
	logger *Logger
)

func InitConfig() {
	logger = NewLogger()
}
func GetLogger() *Logger {
	return logger
}
