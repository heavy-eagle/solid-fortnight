package est

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"

	"go.mozilla.org/pkcs7"
)

func GenerateCSRInMemory(cn string, dns string) ([]byte, crypto.PrivateKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("key generation failed: %v", err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{CommonName: cn},
	}

	if dns != "" {
		template.DNSNames = []string{dns}
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		return nil, nil, fmt.Errorf("CSR creation failed: %v", err)
	}

	return csrDER, privKey, nil
}

func SaveCertAsPEM(p7b []byte, outPath string) error {
	p7, err := pkcs7.Parse(p7b)
	if err != nil {
		return fmt.Errorf("PKCS#7 parse failed: %v", err)
	}

	var buf bytes.Buffer
	for _, cert := range p7.Certificates {
		pem.Encode(&buf, &pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	}

	return os.WriteFile(outPath, buf.Bytes(), 0644)
}

func GenerateCSR(cn, outPath string) error {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate key: %v", err)
	}

	subj := pkix.Name{
		CommonName: cn,
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.ECDSAWithSHA256,
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		return fmt.Errorf("failed to create CSR: %v", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	keyFile := outPath + ".key"
	csrFile := outPath + ".csr"

	keyPEM, err := encodePrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to encode private key: %v", err)
	}

	if err := os.WriteFile(csrFile, csrPEM, 0600); err != nil {
		return fmt.Errorf("failed to write CSR: %v", err)
	}

	if err := os.WriteFile(keyFile, keyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %v", err)
	}

	fmt.Printf("✅ CSR written to %s\n🔐 Private key written to %s\n", csrFile, keyFile)
	return nil
}

func encodePrivateKey(priv crypto.PrivateKey) ([]byte, error) {
	switch k := priv.(type) {
	case *ecdsa.PrivateKey:
		der, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			return nil, err
		}
		return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: der}), nil
	default:
		return nil, fmt.Errorf("unsupported key type")
	}
}
