package wssconfig

type WssConfig struct {
	HostUrl   string `yaml:"hostUrl"`
	Appid     string `yaml:"appid"`
	ApiSecret string `yaml:"apiSecret"`
	ApiKey    string `yaml:"apiKey"`
}
