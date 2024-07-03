package storage

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/minoic/glgf"
	"gopkg.in/yaml.v3"
)

var configOnce sync.Once
var config *ConfigInfo

type ConfigInfo struct {
	InfluxDBToken          string   `yaml:"influxdb_token"`
	InfluxDBURI            string   `yaml:"influxdb_uri"`
	InfluxDBOrg            string   `yaml:"influxdb_org"`
	InfluxDBDatabase       string   `yaml:"influxdb_database"`
	MongoDBURI             string   `yaml:"mongodb_uri"`
	TelegramBotToken       string   `yaml:"telegram_bot_token"`
	TelegramBotApiEndpoint string   `yaml:"telegram_bot_api_endpoint"`
	TelegramBotName        string   `yaml:"telegram_bot_name"`
	SMTPHostname           string   `yaml:"smtp_hostname"`
	SMTPUsername           string   `yaml:"smtp_username"`
	SMTPPassword           string   `yaml:"smtp_password"`
	SMTPPort               int      `yaml:"smtp_port"`
	WebEndpoint            string   `yaml:"web_endpoint"`
	RedisAddrs             []string `yaml:"redis_addrs"`
	RedisUsername          string   `yaml:"redis_username"`
	RedisPassword          string   `yaml:"redis_password"`
	OpenObserveEnabled     bool     `yaml:"openobserve_enabled"`
	OpenObserveOrg         string   `yaml:"openobserve_org"`
	OpenObserveUsername    string   `yaml:"openobserve_username"`
	OpenObservePassword    string   `yaml:"openobserve_password"`
}

func Config() *ConfigInfo {
	configOnce.Do(func() {
		url := os.ExpandEnv("$CONFIG_URL")
		config = new(ConfigInfo)
		resp, err := http.Get(url)
		if err != nil {
			glgf.Error(url)
			panic(err)
		}
		defer resp.Body.Close()
		err = yaml.NewDecoder(resp.Body).Decode(config)
		if err != nil {
			panic(err)
		}
		config.WebEndpoint = strings.TrimRight(config.WebEndpoint, "/")
		config.TelegramBotApiEndpoint = strings.TrimRight(config.TelegramBotApiEndpoint, "/")
	})
	return config
}
