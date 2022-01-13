package xpem

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"testing"

	"github.com/go-pay/gopher/xlog"
)

var (
	pubKey = `-----BEGIN CERTIFICATE-----
MIIDUjCCAjqgAwIBAgIEWRI1ZzANBgkqhkiG9w0BAQsFADA8MR0wGwYDVQQDDBRz
dW5taXNlcnZpY2UgUm9vdCBDQTEOMAwGA1UECgwFc3VubWkxCzAJBgNVBAYTAkNO
MB4XDTIxMTIyOTA2MzAwMFoXDTMxMTEwNzA2MzAwMFowOzEcMBoGA1UEAwwTc3Vu
bWljbGllbnQgUm9vdCBDQTEOMAwGA1UECgwFc3VubWkxCzAJBgNVBAYTAkNOMIIB
IjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnSi2GTKHuJsyP1S2bWQ+Hqug
dotc7WB2qnpJXP7K2LkYEb7ZmW4OIHH+OoQKTlG88+sUAQ/KPiWgRELXtEtG94Mf
7KmAaQcx9Qbek7Z5bh0f2aBsIXRwwVTR9sSMnmX/05kKfMLp1SfF9s4dbzhCtyqj
xpBZyQWEWJFv+/nd/QQibpIg5jzSLhHeYCnHtod86GnxG6gYfKTEdC6vDIj4rOyJ
XGeC7Y4lCevMYInu85VhBxVHvAVN5LIlg8RxEYjVh5VhQ5ghhmKBqSAhqvdTJHXf
8L7Fe5n1XCwfopJMQ/2vPIUlOEGo3WzLsDEBZGTk17IDWsHykZyIDRl9Ux5F/QID
AQABo10wWzAdBgNVHQ4EFgQU+WsDIngFz3R3QhTfWXQ3a/E2rJQwHwYDVR0jBBgw
FoAUM0dERmOJsZw7T4JR+SdOkQPYsbgwDAYDVR0TAQH/BAIwADALBgNVHQ8EBAMC
BeAwDQYJKoZIhvcNAQELBQADggEBAD1o+axYbONV4F37e9MW4b+WRZtxJwqAfUEe
oD905W5FQncLgo2sn3QCv8/linqYNBKzwC538pnhYBuE/CLUsHxI3P47V4qIdW3X
pUg+d3d7/0Pjl8f22GpOwy2DNkk487CUowOQhfmal6xK097gIcAzevjx+Y9/ka40
BMB+HuGGll4b26ozV5Mj2B5vFtmx6MRyplhVEVP5nGztZ6JtHZXQrkphMqbSTPiH
8t4bYygSbUk44gTodre0wx4CYWE17y2Pk45JYzTRAxj85I1O0oGnvcDlFTy2gGis
KVOXz/k2S7zqnlBgV3z6H9Uz6QM4BywwzlgXhEMnPsxJmPbHvQg=
-----END CERTIFICATE-----`
)

func TestDecodeCert(t *testing.T) {
	block, _ := pem.Decode([]byte(pubKey))
	if block == nil {
		return
	}

	pubKeyCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		xlog.Error(err)
		return
	}

	h := sha256.New()
	h.Write(pubKeyCert.Raw)
	s := hex.EncodeToString(h.Sum(nil))
	xlog.Infof("指纹：%s", s)

	//xlog.Infof("%#v", string(pubKeyCert.Raw))
	xlog.Infof("%#v", pubKeyCert.Subject)
	xlog.Infof("序列号：%v", pubKeyCert.SerialNumber)
}
