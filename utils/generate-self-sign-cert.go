package utils

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"

	"github.com/mrzack99s/cocong/vars"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err.Error())
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func GenerateSelfSignCert() {

	_, err := os.Stat("./certs")
	if os.IsNotExist(err) {
		os.MkdirAll("./certs", 0744)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		vars.SystemLog.Println(err.Error())
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"CoCoNG"},
			CommonName:   vars.Config.DomainName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 730),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	template.DNSNames = []string{vars.Config.DomainName}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		fmt.Printf("Failed to create certificate: %s\n", err.Error())
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	err = ioutil.WriteFile("./certs/server.crt", out.Bytes(), 0744)
	if err != nil {
		fmt.Println("Failed to write operator cert file")
	}

	out.Reset()
	pem.Encode(out, pemBlockForKey(priv))

	err = ioutil.WriteFile("./certs/server.key", out.Bytes(), 0744)
	if err != nil {
		fmt.Println("Failed to write operator key file")
	}
	fmt.Println("Generate a self-signed certificate success")
}

func GenerateSelfSignCertWithErrorHandle() error {

	_, err := os.Stat("./certs")
	if os.IsNotExist(err) {
		os.MkdirAll("./certs", 0744)
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"CoCoNG"},
			CommonName:   vars.Config.DomainName,
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 730),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	template.DNSNames = []string{vars.Config.DomainName}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return err
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	err = ioutil.WriteFile("./certs/server.crt", out.Bytes(), 0744)
	if err != nil {
		return err
	}

	out.Reset()
	pem.Encode(out, pemBlockForKey(priv))

	err = ioutil.WriteFile("./certs/server.key", out.Bytes(), 0744)
	if err != nil {
		return err
	}

	return nil
}
