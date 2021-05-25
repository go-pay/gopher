package des

import (
	"testing"

	"github.com/iGoogle-ink/gopher/xlog"
)

var (
	secretKey = "GYBh3Rmey7nNzR/NpV0vAw=="
	iv        = "JR3unO2glQuMhUx3"
)

func TestDesCBCEncryptDecrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCEncryptData([]byte(originData), []byte(secretKey))
	if err != nil {
		xlog.Error("DesCBCEncryptData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCDecryptData(encryptData, []byte(secretKey))
	if err != nil {
		xlog.Error("DesCBCDecryptData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesCBCEncryptDecryptIv(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCEncryptIvData([]byte(originData), []byte(secretKey), []byte(iv))
	if err != nil {
		xlog.Error("DesCBCEncryptIvData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCDecryptIvData(encryptData, []byte(secretKey), []byte(iv))
	if err != nil {
		xlog.Error("DesCBCDecryptIvData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}
