package keycloak

import (
	"crypto/rsa"
	"encoding/hex"
	"math/big"
	"os"
	"testing"
)

func TestPublicKeyParse(t *testing.T) {
	// $ (echo "-----BEGIN PUBLIC KEY-----"; jq -r .public_key internal/keycloak/testdata/publickey.json; echo "-----END PUBLIC KEY-----") | openssl rsa -pubin -inform PEM -text -noout
	// Public-Key: (2048 bit)
	// Modulus:
	//     00:99:86:9c:31:f4:fc:5c:e0:b0:3d:c4:70:66:ed:
	//     b2:d5:c7:d1:91:92:a1:a8:2b:15:e1:8c:68:d5:79:
	//     30:e5:e5:29:14:e6:24:f5:dc:f3:4c:7a:95:1d:54:
	//     5a:b5:22:ef:0b:bf:0b:aa:26:d3:2c:ef:d9:7f:9b:
	//     32:cf:73:8c:b2:03:63:8c:1d:fa:b7:b8:38:4d:73:
	//     03:fe:0e:e8:86:18:07:4e:b6:c5:fc:ad:0a:9d:d9:
	//     aa:cb:d0:6c:cb:76:92:fe:1c:02:a8:9f:a8:54:70:
	//     53:7f:f6:2d:ce:4f:f5:77:02:aa:95:63:57:24:8a:
	//     26:10:28:4a:af:79:dc:89:f6:53:dd:ab:cb:a2:e7:
	//     6d:db:7c:52:81:ec:ca:6d:a8:d5:4f:ad:41:2f:96:
	//     c4:a6:c8:de:1f:32:05:30:e0:e7:ce:d7:ef:02:17:
	//     a3:4e:84:08:20:93:ce:8c:3a:e6:49:75:4b:b5:ab:
	//     5f:22:d6:6e:da:49:07:3c:eb:9b:8f:7c:c2:0a:2f:
	//     62:27:2d:84:31:86:75:fe:69:44:c3:7b:d4:75:80:
	//     b9:f3:99:cc:0b:3b:53:f5:68:18:58:52:59:0e:fb:
	//     98:46:57:b4:23:81:b6:04:fc:4d:26:d3:c6:15:f9:
	//     1e:87:88:49:cc:ed:fb:5b:08:84:f2:8f:45:26:f5:
	//     3b:d1
	// Exponent: 65537 (0x10001)
	nBytes, err := hex.DecodeString(
		`99869c31f4fc5ce0b03dc47066ed` +
			`b2d5c7d19192a1a82b15e18c68d579` +
			`30e5e52914e624f5dcf34c7a951d54` +
			`5ab522ef0bbf0baa26d32cefd97f9b` +
			`32cf738cb203638c1dfab7b8384d73` +
			`03fe0ee88618074eb6c5fcad0a9dd9` +
			`aacbd06ccb7692fe1c02a89fa85470` +
			`537ff62dce4ff57702aa956357248a` +
			`2610284aaf79dc89f653ddabcba2e7` +
			`6ddb7c5281ecca6da8d54fad412f96` +
			`c4a6c8de1f320530e0e7ced7ef0217` +
			`a34e84082093ce8c3ae649754bb5ab` +
			`5f22d66eda49073ceb9b8f7cc20a2f` +
			`62272d84318675fe6944c37bd47580` +
			`b9f399cc0b3b53f568185852590efb` +
			`984657b42381b604fc4d26d3c615f9` +
			`1e878849ccedfb5b0884f28f4526f5` +
			`3bd1`)
	if err != nil {
		t.Fatal(err)
	}
	n := big.NewInt(0)
	n.SetBytes(nBytes)
	var testCases = map[string]struct {
		input  string
		expect *rsa.PublicKey
	}{
		"parse public key JSON": {
			input: "testdata/publickey.json",
			expect: &rsa.PublicKey{
				N: n,
				E: 65537,
			},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(tt *testing.T) {
			f, err := os.Open(tc.input)
			if err != nil {
				tt.Fatal(err)
			}
			pubKey, err := publicKeyParse(f)
			if err != nil {
				tt.Fatal(err)
			}
			if !pubKey.Equal(tc.expect) {
				tt.Fatal("public keys not equal")
			}
		})
	}
}
