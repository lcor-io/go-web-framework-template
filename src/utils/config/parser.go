package config

import "github.com/spf13/viper"

func GetConf() *Config {
	viper.SetConfigName(ConfigFileName)
	viper.SetConfigType(ConfigFileExtension)
	viper.AddConfigPath(ConfigDefaultPath)

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		return nil
	}

	return &c
}
