package config

import "github.com/spf13/viper"

func Read(filepath string, confs ...interface{}) error {
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	for idx := range confs {
		if err := viper.Unmarshal(confs[idx]); err != nil {
			return err
		}
	}
	return nil
}
