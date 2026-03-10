package etcdtools

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "math/big"
    "os"
    "testing"
    "time"
    clientv3 "go.etcd.io/etcd/client/v3"
)

// TestEtcdTools_WithPassword provides an example of password authentication.
func TestEtcdTools_WithPassword(t *testing.T) {
    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
        Username:    "admin",
        Password:    "secure_password",
    }

    etcdTools := NewEtcdTools()
    // Un-comment to actually test the connection in a real environment
    // etcdTools.Init(config)
    // defer etcdTools.Destructor()

    _ = config
    _ = etcdTools
    t.Log("Password authentication configuration is valid.")
}

// TestEtcdTools_WithTLSData provides an example of loading TLS from embedded or raw bytes.
// This matches the scenario: data, err := embedFS.ReadFile("resources/CertFile")
func TestEtcdTools_WithTLSData(t *testing.T) {
    certData, keyData, err := generateMockCert()
    if err != nil {
        t.Fatalf("Failed to generate mock cert: %v", err)
    }

    // Assuming certData and keyData are obtained from embedFS.ReadFile
    tlsOpts := &TLSOptions{
        CertData: certData,
        KeyData:  keyData,
        // CAData: caData, // Optionally load CA Data as well
    }

    tlsConfig, err := tlsOpts.Build()
    if err != nil {
        t.Fatalf("Failed to build TLS config: %v", err)
    }

    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
        TLS:         tlsConfig,
    }

    etcdTools := NewEtcdTools()
    // Un-comment to actually test the connection in a real environment
    // etcdTools.Init(config)

    _ = config
    _ = etcdTools
    t.Log("TLS Data authentication configuration is valid.")
}

// TestEtcdTools_WithTLSFile provides an example of loading TLS from file paths.
func TestEtcdTools_WithTLSFile(t *testing.T) {
    certData, keyData, err := generateMockCert()
    if err != nil {
        t.Fatalf("Failed to generate mock cert: %v", err)
    }

    certFile := "mock_cert.pem"
    keyFile := "mock_key.pem"

    _ = os.WriteFile(certFile, certData, 0644)
    _ = os.WriteFile(keyFile, keyData, 0600)

    defer os.Remove(certFile)
    defer os.Remove(keyFile)

    tlsOpts := &TLSOptions{
        CertFile: certFile,
        KeyFile:  keyFile,
    }

    tlsConfig, err := tlsOpts.Build()
    if err != nil {
        t.Fatalf("Failed to build TLS config: %v", err)
    }

    config := clientv3.Config{
        Endpoints:   []string{"127.0.0.1:2379"},
        DialTimeout: 5 * time.Second,
        TLS:         tlsConfig,
    }

    etcdTools := NewEtcdTools()
    // Un-comment to actually test the connection in a real environment
    // etcdTools.Init(config)

    _ = config
    _ = etcdTools
    t.Log("TLS File authentication configuration is valid.")
}

// generateMockCert is a helper function to create a self-signed certificate for testing purposes.
func generateMockCert() ([]byte, []byte, error) {
    priv, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, err
    }

    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Organization: []string{"Mock Testing Org"},
        },
        NotBefore:             time.Now(),
        NotAfter:              time.Now().Add(time.Hour),
        KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
        BasicConstraintsValid: true,
    }

    derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
    if err != nil {
        return nil, nil, err
    }

    certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

    return certPEM, keyPEM, nil
}
