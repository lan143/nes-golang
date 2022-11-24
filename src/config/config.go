package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	JoyPad1 JoyPadConfig `mapstructure:"joypad1"`
}

func (c *Config) Init() error {
	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			c.initDefaults()
			err := viper.WriteConfigAs("config.yaml")
			if err != nil {
				return errors.Wrap(err, "config.init.read-in-config")
			}
		} else {
			return errors.Wrap(err, "config.init.read-in-config")
		}
	}

	if err := viper.Unmarshal(c); err != nil {
		return errors.Wrap(err, "config.init.unmarshal")
	}

	return nil
}

func (c *Config) initDefaults() {
	log.Println("Create config with default settings...")

	viper.SetDefault("joypad1", JoyPadConfig{
		A:      "Z",
		B:      "X",
		Select: "B",
		Start:  "Enter",
		Up:     "Up",
		Down:   "Down",
		Left:   "Left",
		Right:  "Right",
	})
}

func NewConfig() *Config {
	return &Config{}
}
