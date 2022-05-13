# etcd-tools
可用于服务配置数据的获取和热更新，从而实现依赖配置数据服务的热重载。

## 安装
```shell
go get github.com/lonisy/etcd-tools@v1.0.1
```

## 概述

包 etcdtools 使用了 `go.etcd.io/etcd/client/v3` 客户端。

使用 `etcdtools.NewEtcdTools()` 创建客户端。
```golang
Etcd := etcdtools.NewEtcdTools().Init(clientv3.Config{
    Endpoints:   []string{"xxx.xxx.xxx.xxx:2379","xxx.xxx.xxx.xxx:2379","xxx.xxx.xxx.xxx:2379"},
    DialTimeout: 5 * time.Second,
})
```

使用 viper 来加载服务配置数据。
这里使用 `github.com/spf13/viper` 作为程序应用配置解决方案。
```golang
key := "/projects/apps/domain.com/prod.json"
AppConfig = viper.New()
AppConfig.SetConfigType("json")
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
        // 这里可以触发应用程序的热重载回调
        // Todo ServiceReload... 
    } else {
        log.Panicln(err.Error())
    }
})
```

使用数据结构来加载服务数据
```golang
type Goods []struct {
    GoodsID uint   `json:"id"`
    Name    string `json:"name"`
    Image   string `json:"image"`
}

var GoodsItems Goods

func InitGoodsItems() {
    key := "/projects/apps/domain.com/cache/goods.json"
    Etcd.LoadData(key, func(x, y interface{}) {
        if err := json.Unmarshal([]byte(cast.ToString(y)), &GoodsItems); err != nil {
            log.Panicln(err.Error())
        }
    })
    Etcd.WatchData(key, func(x, y interface{}) {
        if err := json.Unmarshal([]byte(cast.ToString(y)), &GoodsItems); err != nil {
            log.Panicln(err.Error())
        }
    })
}
```