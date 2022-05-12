package example

import (
    clientv3 "go.etcd.io/etcd/client/v3"
    etcdtools "etcd-tools"
    "github.com/spf13/cast"
    "encoding/json"
    "log"
    "time"
)

type Goods []struct {
    GoodsID uint   `json:"id"`
    Name    string `json:"name"`
    Image   string `json:"image"`
}

var GoodsItems Goods

func InitGoodsItems() {
    key := "/projects/apps/domain.com/cache/goods.json"
    Etcd := etcdtools.NewEtcdTools().Init(clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
    })
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
