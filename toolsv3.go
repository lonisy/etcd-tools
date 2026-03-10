package etcdtools

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "errors"
    "fmt"
    "log"
    "os"
    "sync"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
)

// Callback defines the function signature for etcd data events
type Callback func(key, value []byte)

// EtcdTools provides a wrapper around etcd clientv3
type EtcdTools struct {
    once   sync.Once
    Client *clientv3.Client
}

// NewEtcdTools creates a new instance of EtcdTools
func NewEtcdTools() *EtcdTools {
    return &EtcdTools{}
}

// TLSOptions defines the configuration for TLS connection
type TLSOptions struct {
    CertFile   string `json:"cert_file"`
    KeyFile    string `json:"key_file"`
    CAFile     string `json:"ca_file"`
    CertData   []byte `json:"cert_data"`
    KeyData    []byte `json:"key_data"`
    CAData     []byte `json:"ca_data"`
    ServerName string `json:"server_name"`
}

// Build generates a tls.Config based on the provided TLSOptions.
// It supports reading from embedded byte data or from file paths.
func (t *TLSOptions) Build() (*tls.Config, error) {
    var cert tls.Certificate
    var err error

    // Load Certificate and Key
    if len(t.CertData) > 0 && len(t.KeyData) > 0 {
        cert, err = tls.X509KeyPair(t.CertData, t.KeyData)
    } else if t.CertFile != "" && t.KeyFile != "" {
        cert, err = tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
    } else {
        return nil, errors.New("missing certificate or key data/file")
    }

    if err != nil {
        return nil, fmt.Errorf("failed to load client certificate: %w", err)
    }

    tlsConfig := &tls.Config{
        Certificates:       []tls.Certificate{cert},
        ServerName:         t.ServerName,
        InsecureSkipVerify: true,
    }

    // Load CA if provided
    var caCertPool *x509.CertPool
    if len(t.CAData) > 0 {
        caCertPool = x509.NewCertPool()
        if !caCertPool.AppendCertsFromPEM(t.CAData) {
            return nil, errors.New("failed to parse CA certificate from data")
        }
        tlsConfig.RootCAs = caCertPool
    } else if t.CAFile != "" {
        caData, err := os.ReadFile(t.CAFile)
        if err != nil {
            return nil, fmt.Errorf("failed to read CA file: %w", err)
        }
        caCertPool = x509.NewCertPool()
        if !caCertPool.AppendCertsFromPEM(caData) {
            return nil, errors.New("failed to parse CA certificate from file")
        }
        tlsConfig.RootCAs = caCertPool
    }

    return tlsConfig, nil
}

// Init initializes the etcd client with the provided configuration
func (s *EtcdTools) Init(config clientv3.Config) *EtcdTools {
    s.once.Do(func() {
        var err error
        s.Client, err = clientv3.New(config)
        if err != nil {
            log.Panicf("failed to create etcd client: %v", err)
        }
    })
    return s
}

// LoadData retrieves a key from etcd and executes the callback
func (s *EtcdTools) LoadData(key string, callback Callback) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
    defer cancel()

    resp, err := s.Client.Get(ctx, key)
    if err != nil {
        log.Panicf("failed to get key '%s': %v", key, err)
    }

    if resp.Count > 0 {
        for _, ev := range resp.Kvs {
            if string(ev.Key) == key {
                callback(ev.Key, ev.Value)
            }
        }
    } else {
        log.Panicf("please check etcd key '%s', not found", key)
    }
}

// WatchData watches a specific key and executes the callback upon events
func (s *EtcdTools) WatchData(key string, callback Callback) {
    go func(cli *clientv3.Client) {
        log.Println("etcd watching...")
        rch := cli.Watch(context.Background(), key)
        for wresp := range rch {
            if wresp.Canceled {
                log.Printf("watch canceled for key: %s\n", key)
                return
            }
            if wresp.Err() != nil {
                log.Printf("watch error for key %s: %v\n", key, wresp.Err())
                return
            }
            for _, ev := range wresp.Events {
                log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
                if string(ev.Kv.Key) == key {
                    callback(ev.Kv.Key, ev.Kv.Value)
                }
            }
        }
    }(s.Client)
}

// Destructor gracefully closes the etcd client connection
func (s *EtcdTools) Destructor() {
    if s.Client != nil {
        if err := s.Client.Close(); err != nil {
            log.Printf("failed to close etcd client: %v\n", err)
        }
    }
}
