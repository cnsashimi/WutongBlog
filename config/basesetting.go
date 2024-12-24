package config

func Getsetting() Base {
	configfile := GetYml()
	return configfile.Base
}
