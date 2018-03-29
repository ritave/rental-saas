package config

import (
	"github.com/spf13/viper"
	"log"
)

type C struct {
	CORS     []string `json:"CORS"`
	Callback string   `json:"callback:"`
	Receiver struct {
		Expiration int    `json:"expiration"`
		Channel    string `json:"channel"`
	} `json:"receiver"`
	Db struct {
		Restart bool `json:"restart"`
	} `json:"db"`
	Server struct {
		Self    string `json:"self"`
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"server"`
}

func GetConfig() (C) {
	v := viper.New()
	v.AddConfigPath(".")
	v.AddConfigPath("secrets")
	v.SetConfigName("config")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read in config: %s", err.Error())
	}

	c := C{}

	v.Unmarshal(&c)
	return c
}
