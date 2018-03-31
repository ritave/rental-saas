package config

import (
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

type C struct {
	CORS     []string `json:"CORS"`
	Callback string   `json:"callback:"`
	Receiver struct {
		Expiration int    `json:"expiration"`
		Channel    string `json:"channel"`
	} `json:"receiver"`
	DB struct {
		Restart bool `json:"restart"`
	} `json:"db"`
	Server struct {
		Self    string `json:"self"`
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"server"`
	Logging struct {
		Debug         bool   `json:"debug"`
		StoreInDB     bool   `json:"storeInDB"`
		SomethingElse string `json:"somethingElse"`
	} `json:"logging"`
}

func GetConfig() (C) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("secrets")
	v.SetConfigName("config")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Failed to read in config: %s", err.Error())
	}

	c := C{}
	v.Unmarshal(&c)
	return c
}
