package infra

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

const (
	keySize        = 3072
	keyExpiryYears = 2
)

var CertSubject = pkix.Name{
	CommonName:         "Dependabot Internal CA - Lite",
	OrganizationalUnit: []string{"Dependabot Lite"},
	Organization:       []string{"GitHub Inc."},
	Locality:           []string{"San Francisco"},
	Province:           []string{"California"},
	Country:            []string{"US"},
}

func GenerateCertificateAuthority() (CertificateAuthority, error) {
	key, pemKey, err := generateKey()
	if err != nil {
		return CertificateAuthority{}, err
	}
	pemCert, err := generateCert(key)
	if err != nil {
		return CertificateAuthority{}, err
	}
	return CertificateAuthority{Cert: pemCert, Key: pemKey}, nil
}

func generateKey() (*rsa.PrivateKey, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, "", err
	}
	kb := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return key, string(pem.EncodeToMemory(kb)), nil
}

func generateCert(key *rsa.PrivateKey) (string, error) {
	serial, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.AddDate(keyExpiryYears, 0, 0)

	pubBytes, _ := x509.MarshalPKIXPublicKey(key.Public())
	h := sha1.Sum(pubBytes)
	subjectKeyId := h[:]

	template := x509.Certificate{
		SerialNumber:          serial,
		Subject:               CertSubject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		SubjectKeyId:          subjectKeyId,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, key.Public(), key)
	if err != nil {
		return "", err
	}
	cb := &pem.Block{Type: "CERTIFICATE", Bytes: cert}
	return string(pem.EncodeToMemory(cb)), nil
}
