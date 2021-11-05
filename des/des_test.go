package des

import (
	"testing"

	"github.com/go-pay/gopher/xlog"
)

var (
	secretKey1 = "hehehaha"
	iv1        = "12378945"

	secretKey = "GYBh3Rmey7nNzR/NpV0vAw=="
	iv        = "JR3unO2glQuMhUx3"
)

func TestDesCBCEncryptDecrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCEncryptData([]byte(originData), []byte(secretKey1))
	if err != nil {
		xlog.Error("DesCBCEncryptData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCDecryptData(encryptData, []byte(secretKey1))
	if err != nil {
		xlog.Error("DesCBCDecryptData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesCBCEncryptDecryptIv(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCEncryptIvData([]byte(originData), []byte(secretKey1), []byte(iv1))
	if err != nil {
		xlog.Error("DesCBCEncryptIvData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCDecryptIvData(encryptData, []byte(secretKey1), []byte(iv1))
	if err != nil {
		xlog.Error("DesCBCDecryptIvData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesCBCTripleEncryptDecrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCTripleEncryptData([]byte(originData), []byte(secretKey))
	if err != nil {
		xlog.Error("DesCBCEncryptData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCTripleDecryptData(encryptData, []byte(secretKey))
	if err != nil {
		xlog.Error("DesCBCDecryptData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesCBCTripleEncryptDecryptIv(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCTripleEncryptIvData([]byte(originData), []byte(secretKey), []byte(iv))
	if err != nil {
		xlog.Error("DesCBCEncryptIvData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCTripleDecryptIvData(encryptData, []byte(secretKey), []byte(iv))
	if err != nil {
		xlog.Error("DesCBCDecryptIvData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}
