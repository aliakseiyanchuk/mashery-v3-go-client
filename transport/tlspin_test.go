package transport_test

import (
	"context"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/transport"
	"github.com/aliakseiyanchuk/mashery-v3-go-client/v3client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultPinnerWillAcceptMashery(t *testing.T) {
	wildCl := v3client.NewWildcardClient(v3client.Params{
		MashEndpoint: "https://api.mashery.com",
	})

	_, err := wildCl.FetchAny(context.Background(), "", nil)
	assert.Nil(t, err)
}

func TestDefaultPinnerWillRejectGoogle(t *testing.T) {
	wildCl := v3client.NewWildcardClient(v3client.Params{
		MashEndpoint: "https://www.google.com",
	})

	_, err := wildCl.FetchAny(context.Background(), "", nil)
	assert.NotNil(t, err)
}

//
//func TestPinnerWillSuccessfullyValidate(t *testing.T) {
//	p := &transport.TLSPinner{}
//
//	leafPin := transport.TLSCertChainPin{
//		CommonName: "*.mashery.com",
//	}
//	leafPin.SerialNumberFromHex("0165f717537bc5454528a1147389b010")
//	leafPin.FingerprintFrom("6592068a059c5ae724f6c8ad23048799c87b47b0143d846eb603ad5f09597473")
//
//	p.Add(leafPin)
//
//	issuerPin := transport.TLSCertChainPin{
//		CommonName: "DigiCert Global G2 TLS RSA SHA256 2020 CA1",
//	}
//	issuerPin.SerialNumberFromHex("0cf5bd062b5602f47ab8502c23ccf066")
//	issuerPin.FingerprintFrom("c8025f9fc65fdfc95b3ca8cc7867b9a587b5277973957917463fc813d0b625a9")
//	p.Add(issuerPin)
//
//	wildCl := v3client.NewWildcardClient(v3client.Params{
//		MashEndpoint: "https://api.mashery.com",
//		HTTPClientParams: transport.HTTPClientParams{
//			TLSConfig: p.CreateTLSConfig(),
//		},
//	})
//
//	_, err := wildCl.FetchAny(context.Background(), "", nil)
//	assert.Nil(t, err)
//
//}

func TestPinnerWillSuccessfullyRejectCertOnMismatch(t *testing.T) {
	p := &transport.TLSPinner{}

	leafPin := transport.TLSCertChainPin{
		CommonName: "*.mashery.com",
	}
	leafPin.SerialNumberFromHex("0D:DB:59:B2:0C:1D:CD:CF:3A:CC:E1:74:90:70:78:3D")
	leafPin.FingerprintFrom("42:7A:76:D0:92:4E:DB:29:A9:8F:92:4A:BA:C4:5C:71:5D:63:75:7C:FE:72:6B:B1:88:20:A9:C5:61:52:AA:8D")

	p.Add(leafPin)

	issuerPin := transport.TLSCertChainPin{
		CommonName: "DigiCert TLS RSA SHA256 2020 CA1",
	}
	issuerPin.SerialNumberFromHex("0A:35:08:D5:5C:29:2B:01:7D:F8:AD:65:C0:0F:F7:E4")
	issuerPin.FingerprintFrom("26:76:87:13:D3:B4:59:F9:38:2D:2A:59:4F:85:F3:47:09:FD:2A:89:30:73:15:42:A4:14:6F:FB:24:6B:EC:69")
	p.Add(issuerPin)

	wildCl := v3client.NewWildcardClient(v3client.Params{
		MashEndpoint: "https://api.mashery.com",
		HTTPClientParams: transport.HTTPClientParams{
			TLSConfig: p.CreateTLSConfig(),
		},
	})

	_, err := wildCl.FetchAny(context.Background(), "", nil)
	assert.NotNil(t, err)

}
