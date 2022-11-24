package config

type JoyPadConfig struct {
	A      string `mapstructure:"a"`
	B      string `mapstructure:"b"`
	Select string `mapstructure:"select"`
	Start  string `mapstructure:"start"`
	Up     string `mapstructure:"up"`
	Down   string `mapstructure:"down"`
	Left   string `mapstructure:"left"`
	Right  string `mapstructure:"right"`
}
