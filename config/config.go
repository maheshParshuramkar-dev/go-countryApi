package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"

	"github.com/go-countryApi/models"
)

var AppConfig *models.Configurations

func GetConfig() *models.Configurations {
	if AppConfig != nil {
		return AppConfig
	}

	// get current directory path
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)
	var configPath string

	// name of config file (without extension)
	viper.SetConfigName("config")

	// where to fetch config file from
	if os.Getenv("GO_ENV") != "local" {
		configPath = "/etc/config"
	} else {
		configPath = basepath + "/data"
	}

	// add path where to search for config file
	viper.AddConfigPath(configPath)

	// set the file type
	viper.SetConfigType("json")

	// read config from the specified path
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// unmarshal read configurations into a struct
	err := viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error while decoding config file: %s", err))
	}

	return AppConfig
}
