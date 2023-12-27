package crypto

import (
	"crypto/x509"
	"fmt"
	"time"
)

type CertVerifier func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error

func PeerCertificateVerifier(caCertPEM []byte) CertVerifier {
	return func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		certs := make([]*x509.Certificate, len(rawCerts))
		for i, asn1Data := range rawCerts {
			cert, err := x509.ParseCertificate(asn1Data)
			if err != nil {
				return fmt.Errorf("failed to parse certificate: %v", err)
			}
			certs[i] = cert
		}
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCertPEM)
		opts := x509.VerifyOptions{
			Roots:         certPool,
			CurrentTime:   time.Now(),
			DNSName:       "", // Skip hostname verification
			Intermediates: x509.NewCertPool(),
		}

		for i, cert := range certs {
			if i == 0 {
				continue
			}
			opts.Intermediates.AddCert(cert)
		}
		_, err := certs[0].Verify(opts)
		return err
	}
}
