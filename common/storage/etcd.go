package storage

import (
	"github.com/minoic/glgf"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"strings"
	"sync"
)

var etcdOnce sync.Once
var etcdConfig *ConfigInfo

type ConfigInfo struct {
	InfluxDBToken          string   `yaml:"influxdb_token"`
	InfluxDBURI            string   `yaml:"influxdb_uri"`
	InfluxDBOrg            string   `yaml:"influxdb_org"`
	InfluxDBDatabase       string   `yaml:"influxdb_database"`
	MongoDBURI             string   `yaml:"mongodb_uri"`
	TelegramBotToken       string   `yaml:"telegram_bot_token"`
	TelegramBotApiEndpoint string   `yaml:"telegram_bot_api_endpoint"`
	TelegramBotName        string   `yaml:"telegram_bot_name"`
	TencentSTMPPassword    string   `yaml:"tencent_smtp_password"`
	WebEndpoint            string   `yaml:"web_endpoint"`
	RedisAddrs             []string `yaml:"redis_addrs"`
	RedisUsername          string   `yaml:"redis_username"`
	RedisPassword          string   `yaml:"redis_password"`
}

func Config() *ConfigInfo {
	etcdOnce.Do(func() {
		url := os.ExpandEnv("$CONFIG_URL")
		etcdConfig = new(ConfigInfo)
		resp, err := http.Get(url)
		if err != nil {
			glgf.Error(url)
			panic(err)
		}
		defer resp.Body.Close()
		err = yaml.NewDecoder(resp.Body).Decode(etcdConfig)
		if err != nil {
			panic(err)
		}
		etcdConfig.WebEndpoint = strings.TrimRight(etcdConfig.WebEndpoint, "/")
		etcdConfig.TelegramBotApiEndpoint = strings.TrimRight(etcdConfig.TelegramBotApiEndpoint, "/")
	})
	return etcdConfig
}
