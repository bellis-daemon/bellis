package storage

import (
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"sync"
)

var etcdOnce sync.Once
var etcdConfig *viper.Viper

func Config() *viper.Viper {
	etcdOnce.Do(func() {
		etcdConfig = viper.New()
		err := etcdConfig.AddRemoteProvider("etcd3", "etcd:2379", "/config")
		if err != nil {
			panic(err)
		}
		etcdConfig.SetConfigType("yaml") // Need to explicitly set this to json
		err = etcdConfig.ReadRemoteConfig()
		if err != nil {
			panic(err)
		}
	})
	return etcdConfig
}
