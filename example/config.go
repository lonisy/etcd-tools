package example

import (
    "bytes"
    etcdtools "etcd-tools"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cast"
    "github.com/spf13/viper"
    clientv3 "go.etcd.io/etcd/client/v3"
    "time"
)

var AppConfig *viper.Viper

func InitConfig() {
    key := "/projects/apps/domain.com/prod.json"
    AppConfig = viper.New()
    AppConfig.SetConfigType("json")
    Etcd := etcdtools.NewEtcdTools().Init(clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
    })
    Etcd.LoadData(key, func(x, y interface{}) {
        if err := AppConfig.ReadConfig(bytes.NewBufferString(cast.ToString(y))); err == nil {
            if err := AppConfig.WriteConfigAs("config.json"); err != nil {
                log.Println(err.Error())
            }
        } else {
            log.Panicln(err.Error())
        }
    })
    Etcd.WatchData(key, func(x, y interface{}) {
        if err := AppConfig.ReadConfig(bytes.NewBufferString(cast.ToString(y))); err == nil {
            if err := AppConfig.WriteConfigAs("config.json"); err != nil {
                log.Println(err.Error())
            }
            // Todo ServiceReload...
        } else {
            log.Panicln(err.Error())
        }
    })
}
