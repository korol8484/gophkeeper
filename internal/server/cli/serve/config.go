package serve

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"github.com/korol8484/gophkeeper/internal/server/app"
	"github.com/korol8484/gophkeeper/internal/server/db"
	"github.com/korol8484/gophkeeper/internal/server/token"
	"github.com/spf13/viper"
	"math/big"
	"net"
	"os"
	"path"
	"time"
)

type Config struct {
	Db    *db.Config    `mapstructure:"db"`
	Token *token.Config `mapstructure:"token"`
	Http  *app.Config   `mapstructure:"http"`
}

func newConfig() (*Config, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can't retrive pwd %w", err)
	}

	cfg := &Config{
		Db: &db.Config{
			MaxIdleConn:     1,
			MaxOpenConn:     10,
			MaxLifetimeConn: time.Minute * 3,
		},
		Token: &token.Config{
			Secret: "1234567891aaa",
			Name:   "Authorization",
			Expire: 24 * time.Hour,
		},
		Http: &app.Config{
			Listen:         ":8199",
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Key:            path.Join(pwd, "/server.key"),
			Pem:            path.Join(pwd, "/server.pem"),
		},
	}

	if err = viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("can't unmarshal config: %w", err)
	}

	err = createTLS(cfg.Http.Pem, cfg.Http.Key)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func createTLS(pemPath string, keyPath string) error {
	cert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Shortener"},
			Country:      []string{"RU"},
			Province:     []string{"Moscow"},
			Locality:     []string{"Moscow"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	certBytes, _ := x509.CreateCertificate(rand.Reader, &cert, &cert, &privateKey.PublicKey, privateKey)
	err := saveCertToFile(pemPath, "CERTIFICATE", certBytes)
	if err != nil {
		return err
	}

	err = saveCertToFile(keyPath, "RSA PRIVATE KEY", x509.MarshalPKCS1PrivateKey(privateKey))
	if err != nil {
		return err
	}

	return nil
}

func saveCertToFile(filePath string, cypherType string, cypher []byte) error {
	var (
		buf  bytes.Buffer
		file *os.File
	)

	err := pem.Encode(&buf, &pem.Block{
		Type:  cypherType,
		Bytes: cypher,
	})
	if err != nil {
		return fmt.Errorf("can't encode pem: %w", err)
	}

	file, _ = os.Create(filePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	_, err = buf.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}
