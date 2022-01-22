package transport

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"
)

type TLSCertChainPin struct {
	CommonName   string
	SerialNumber []byte
	Fingerprint  []byte
}

func (pin *TLSCertChainPin) SerialNumberFromHex(str string) error {
	return decodeUserHex(str, &pin.SerialNumber)
}

func (pin *TLSCertChainPin) FingerprintFrom(str string) error {
	return decodeUserHex(str, &pin.Fingerprint)
}

func (pin *TLSCertChainPin) IsEmpty() bool {
	return len(pin.CommonName) == 0 && len(pin.SerialNumber) == 0 && len(pin.Fingerprint) == 0
}

func decodeUserHex(input string, output *[]byte) error {
	normalized := strings.ReplaceAll(input, ":", "")
	decoded, err := hex.DecodeString(normalized)
	if err == nil {
		*output = decoded
	}

	return err
}

func (pin TLSCertChainPin) PinnedSerial() *big.Int {
	if len(pin.SerialNumber) > 0 {
		v := big.Int{}
		v.SetBytes(pin.SerialNumber)

		return &v
	} else {
		return nil
	}
}

type TLSPinner struct {
	TLSCertChainPins []TLSCertChainPin
}

func (pinner *TLSPinner) Add(pin TLSCertChainPin) *TLSPinner {
	pinner.TLSCertChainPins = append(pinner.TLSCertChainPins, pin)
	return pinner
}

func (tlp *TLSPinner) verifyPeerCertificateImpl(_ [][]byte, verifiedChains [][]*x509.Certificate) error {
chains:
	for _, v := range verifiedChains {

	pins:
		for _, chainPin := range tlp.TLSCertChainPins {

			pinnedSerial := chainPin.PinnedSerial()
			checkTime := time.Now()

			for _, cert := range v {
				// We can't accept expired certificates in the chain path
				if checkTime.Before(cert.NotBefore) || checkTime.After(cert.NotAfter) {
					continue chains
				}

				// Check the common name pin
				if len(chainPin.CommonName) > 0 && chainPin.CommonName != cert.Subject.CommonName {
					continue
				}

				if pinnedSerial != nil {
					if cert.SerialNumber == nil {
						continue
					} else if cert.SerialNumber.Cmp(pinnedSerial) != 0 {
						continue
					}
				}

				if len(chainPin.Fingerprint) > 0 {
					hash := sha256.New()
					hash.Write(cert.Raw)

					sig := hash.Sum(nil)
					if bytes.Compare(chainPin.Fingerprint, sig) != 0 {
						continue
					}
				}

				// Found a pinned certificate this chain. Let's look for
				// next PIN
				continue pins
			}

			// We have a pin that does not occur in this chain.
			// Let's continue searching for next chain
			continue chains
		}

		// This chain contains all required pins.
		return nil
	}

	return errors.New("no matching chains")
}

func (tlp *TLSPinner) CreateTLSConfig() *tls.Config {
	return &tls.Config{
		VerifyPeerCertificate: tlp.verifyPeerCertificateImpl,
	}
}

var defaultMasheryPinner *TLSPinner

func init() {
	defaultMasheryPinner = &TLSPinner{}

	leafPIN := TLSCertChainPin{
		CommonName: "*.mashery.com",
	}
	defaultMasheryPinner.Add(leafPIN)

	issuerPin := TLSCertChainPin{
		CommonName: "DigiCert TLS RSA SHA256 2020 CA1",
	}
	issuerPin.SerialNumberFromHex("0A:35:08:D5:5C:29:2B:01:7D:F8:AD:65:C0:0F:F7:E4")
	issuerPin.FingerprintFrom("25:76:87:13:D3:B4:59:F9:38:2D:2A:59:4F:85:F3:47:09:FD:2A:89:30:73:15:42:A4:14:6F:FB:24:6B:EC:69")
	defaultMasheryPinner.Add(issuerPin)
}

func DefaultTLSConfig() *tls.Config {
	return defaultMasheryPinner.CreateTLSConfig()
}
