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
		CommonName: "DigiCert Global G2 TLS RSA SHA256 2020 CA1",
	}
	issuerPin.SerialNumberFromHex("0C:F5:BD:06:2B:56:02:F4:7A:B8:50:2C:23:CC:F0:66")
	issuerPin.FingerprintFrom("C8:02:5F:9F:C6:5F:DF:C9:5B:3C:A8:CC:78:67:B9:A5:87:B5:27:79:73:95:79:17:46:3F:C8:13:D0:B6:25:A9")
	defaultMasheryPinner.Add(issuerPin)
}

func DefaultTLSConfig() *tls.Config {
	return defaultMasheryPinner.CreateTLSConfig()
}
