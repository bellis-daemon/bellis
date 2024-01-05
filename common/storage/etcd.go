package storage

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
	"sync"
)

var etcdOnce sync.Once
var etcdConfig *ConfigInfo

type ConfigInfo struct {
	InfluxDBToken          string `mapstructure:"influxdb_token"`
	MongoDBURI             string `mapstructure:"mongodb_uri"`
	TelegramBotToken       string `mapstructure:"telegram_bot_token"`
	TelegramBotApiEndpoint string `mapstructure:"telegram_bot_api_endpoint"`
	TelegramBotName        string `mapstructure:"telegram_bot_name"`
	TencentSTMPPassword    string `mapstructure:"tencent_smtp_password"`
	WebEndpoint            string `mapstructure:"web_endpoint"`
}

func Config() *ConfigInfo {
	etcdOnce.Do(func() {
		etcdConfig = new(ConfigInfo)
		err := viper.AddRemoteProvider("etcd3", "etcd:2379", "/config")
		if err != nil {
			panic(err)
		}
		viper.SetConfigType("yaml") // Need to explicitly set this to json
		err = viper.ReadRemoteConfig()
		if err != nil {
			panic(err)
		}
		err = viper.Unmarshal(etcdConfig)
		if err != nil {
			panic(err)
		}
		etcdConfig.WebEndpoint = strings.TrimRight(etcdConfig.WebEndpoint, "/")
		etcdConfig.TelegramBotApiEndpoint = strings.TrimRight(etcdConfig.TelegramBotApiEndpoint, "/")
	})
	return etcdConfig
}
