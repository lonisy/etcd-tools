package etcdtools

import (
    "context"
    "fmt"
    clientv3 "go.etcd.io/etcd/client/v3"
    "log"
    "sync"
    "time"
)

type Callback func(x, y interface{})

type EtcdTools struct {
    once   sync.Once
    Client *clientv3.Client
}

func NewEtcdTools() *EtcdTools {
    return &EtcdTools{}
}

func (s *EtcdTools) Init(config clientv3.Config) *EtcdTools {
    s.once.Do(func() {
        var err error
        s.Client, err = clientv3.New(config)
        if err != nil {
            log.Panic(err.Error())
        }
    })
    return s
}

func (s *EtcdTools) LoadData(key string, callback Callback) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
    defer cancel()
    if resp, err := s.Client.Get(ctx, key); err == nil {
        if resp.Count > 0 {
            for _, ev := range resp.Kvs {
                if string(ev.Key) == key {
                    callback(ev.Key, ev.Value)
                }
            }
        } else {
            log.Panic(fmt.Sprintf("please check etcd key '%?'", key))
        }
    } else {
        log.Panic(err.Error())
    }
}

func (s *EtcdTools) WatchData(key string, callback Callback) {
    go func(cli *clientv3.Client) {
        log.Println("etcd watching...")
        rch := cli.Watch(context.Background(), key)
        for wresp := range rch {
            for _, ev := range wresp.Events {
                log.Println(fmt.Sprintf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value))
                if string(ev.Kv.Key) == key {
                    callback(ev.Kv.Key, ev.Kv.Value)
                }
            }
        }
    }(s.Client)
}

func (s *EtcdTools) Destructor() {
    if err := s.Client.Close(); err != nil {
        log.Println(err.Error())
    }
}
