package config

import (
	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

type C struct {
	CORS     []string `json:"CORS"`
	Google   struct {
		Receiver struct {
			Expiration int    `json:"expiration"`
			Channel    string `json:"channel"`
		} `json:"receiver"`
	} `json:"google"`
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
	Pozamiatane struct {
		Address string `json:"address"`
	} `json:"pozamiatane"`
}

func GetConfig() (C) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("configs")
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
