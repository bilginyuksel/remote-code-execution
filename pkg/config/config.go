package config

import "github.com/spf13/viper"

func Read(filepath string, conf interface{}) error {
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(conf)
}
