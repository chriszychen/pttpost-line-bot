package config

import "github.com/jinzhu/configor"

var Config = struct {
	ChannelSecret      string `env:"CHANNEL_SECRET" yaml:"ChannelSecret" default:""`
	ChannelAccessToken string `env:"CHANNEL_ACCESS_TOKEN" yaml:"ChannelAccessToken" default:""`
	SelfLineID         string `env:"SELF_LINE_ID" yaml:"SelfLineID" default:""`
}{}

func Init() {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		panic(err)
	}
}
