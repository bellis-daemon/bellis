package models

import (
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/youmark/pkcs8"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sort"
	"strings"
)

const TLSMinVersionDefault = tls.VersionTLS12

type TLS struct {
	ID     primitive.ObjectID `json:"ID" bson:"_id"`
	UserID primitive.ObjectID `json:"UserID" bson:"UserID"`
	Name   string             `json:"Name" bson:"Name"`
	// TLS Settings
	TLSCA         string `json:"TLSCA" bson:"TLSCA"`
	TLSCert       string `json:"TLSCert" bson:"TLSCert"`
	TLSKey        string `json:"TLSKey" bson:"TLSKey"`
	TLSKeyPwd     string `json:"TLSKeyPwd" bson:"TLSKeyPwd"`
	TLSMinVersion string `json:"TLSMinVersion" bson:"TLSMinVersion"` // "TLS10" "TLS11" "TLS12" "TLS13"
	Insecure      bool   `json:"Insecure" bson:"Insecure"`
}

func (this *TLS) TLSConfig() (*tls.Config, error) {
	empty := this.TLSCA == "" && this.TLSKey == "" && this.TLSCert == ""

	if empty {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: this.Insecure,
		Renegotiation:      tls.RenegotiateNever,
	}

	if this.TLSCA != "" {
		pool, err := makeCertPool([]string{this.TLSCA})
		if err != nil {
			return nil, err
		}
		tlsConfig.RootCAs = pool
	}

	if this.TLSCert != "" && this.TLSKey != "" {
		err := loadCertificate(tlsConfig, this.TLSCert, this.TLSKey, this.TLSKeyPwd)
		if err != nil {
			return nil, err
		}
	}

	// Explicitly and consistently set the minimal accepted version using the
	// defined default. We use this setting for both clients and servers
	// instead of relying on Golang's default that is different for clients
	// and servers and might change over time.
	tlsConfig.MinVersion = TLSMinVersionDefault
	if this.TLSMinVersion != "" {
		version, err := ParseTLSVersion(this.TLSMinVersion)
		if err != nil {
			return nil, fmt.Errorf("could not parse tls min version %q: %w", this.TLSMinVersion, err)
		}
		tlsConfig.MinVersion = version
	}

	return tlsConfig, nil
}

func loadCertificate(config *tls.Config, certString, keyString, privateKeyPassphrase string) error {
	var err error
	certBytes := []byte(certString)

	keyBytes := []byte(keyString)

	keyPEMBlock, _ := pem.Decode(keyBytes)
	if keyPEMBlock == nil {
		return errors.New("failed to decode private key: no PEM data found")
	}

	var cert tls.Certificate
	if keyPEMBlock.Type == "ENCRYPTED PRIVATE KEY" {
		if privateKeyPassphrase == "" {
			return errors.New("missing password for PKCS#8 encrypted private key")
		}
		var decryptedKey *rsa.PrivateKey
		decryptedKey, err = pkcs8.ParsePKCS8PrivateKeyRSA(keyPEMBlock.Bytes, []byte(privateKeyPassphrase))
		if err != nil {
			return fmt.Errorf("failed to parse encrypted PKCS#8 private key: %w", err)
		}
		cert, err = tls.X509KeyPair(certBytes, pem.EncodeToMemory(&pem.Block{Type: keyPEMBlock.Type, Bytes: x509.MarshalPKCS1PrivateKey(decryptedKey)}))
		if err != nil {
			return fmt.Errorf("failed to load cert/key pair: %w", err)
		}
	} else if keyPEMBlock.Headers["Proc-Type"] == "4,ENCRYPTED" {
		// The key is an encrypted private key with the DEK-Info header.
		// This is currently unsupported because of the deprecation of x509.IsEncryptedPEMBlock and x509.DecryptPEMBlock.
		return fmt.Errorf("password-protected keys in pkcs#1 format are not supported")
	} else {
		cert, err = tls.X509KeyPair(certBytes, keyBytes)
		if err != nil {
			return fmt.Errorf("failed to load cert/key pair: %w", err)
		}
	}
	config.Certificates = []tls.Certificate{cert}
	return nil
}

var tlsVersionMap = map[string]uint16{
	"TLS10": tls.VersionTLS10,
	"TLS11": tls.VersionTLS11,
	"TLS12": tls.VersionTLS12,
	"TLS13": tls.VersionTLS13,
}

// ParseTLSVersion returns a `uint16` by received version string key that represents tls version from crypto/tls.
// If version isn't supported ParseTLSVersion returns 0 with error
func ParseTLSVersion(version string) (uint16, error) {
	if v, ok := tlsVersionMap[version]; ok {
		return v, nil
	}

	available := make([]string, 0, len(tlsVersionMap))
	for n := range tlsVersionMap {
		available = append(available, n)
	}
	sort.Strings(available)
	return 0, fmt.Errorf("unsupported version %q (available: %s)", version, strings.Join(available, ","))
}

func makeCertPool(certs []string) (*x509.CertPool, error) {
	pool := x509.NewCertPool()
	for _, cert := range certs {
		if !pool.AppendCertsFromPEM([]byte(cert)) {
			return nil, fmt.Errorf("could not parse any PEM certificates %s", cert)
		}
	}
	return pool, nil
}
