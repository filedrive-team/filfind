package settings

import (
	"fmt"
	"github.com/jinzhu/configor"
)

func LoadConfig(configFile string) *AppConfig {
	tomlConfig := new(AppConfig)
	err := configor.Load(tomlConfig, configFile)
	if err != nil {
		panic(fmt.Sprintf("fail to load app config:\n %v\n", err))
	}
	return tomlConfig
}
