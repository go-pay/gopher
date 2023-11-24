package des

import (
	"testing"

	"github.com/go-pay/xlog"
)

var (
	secretKey1 = "hehehaha"
	iv1        = "12378945"

	secretKey = "GYBh3Rmey7nNzR/NpV0vAw=="
	iv        = "JR3unO2glQuMhUx3"
)

func init() {
	xlog.Level = xlog.DebugLevel
}

func TestDesCBCEncrypt_Decrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCEncrypt([]byte(originData), []byte(secretKey1))
	if err != nil {
		xlog.Error("DesCBCEncryptData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCDecrypt(encryptData, []byte(secretKey1))
	if err != nil {
		xlog.Error("DesCBCDecryptData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesECBEncrypt_Decrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := ECBEncrypt([]byte(originData), []byte(secretKey1), []byte(iv1))
	if err != nil {
		xlog.Error("DesCBCEncryptIvData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := ECBDecrypt(encryptData, []byte(secretKey1), []byte(iv1))
	if err != nil {
		xlog.Error("DesCBCDecryptIvData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesCBCTripleEncrypt_Decrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := CBCTripleEncrypt([]byte(originData), []byte(secretKey))
	if err != nil {
		xlog.Error("DesCBCEncryptData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := CBCTripleDecrypt(encryptData, []byte(secretKey))
	if err != nil {
		xlog.Error("DesCBCDecryptData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}

func TestDesECBTripleEncrypt_Decrypt(t *testing.T) {
	originData := "www.gopay.ink"
	xlog.Debug("originData:", originData)
	encryptData, err := ECBTripleEncrypt([]byte(originData), []byte(secretKey), []byte(iv))
	if err != nil {
		xlog.Error("DesCBCEncryptIvData:", err)
		return
	}
	xlog.Debug("encryptData:", string(encryptData))
	origin, err := ECBTripleDecrypt(encryptData, []byte(secretKey), []byte(iv))
	if err != nil {
		xlog.Error("DesCBCDecryptIvData:", err)
		return
	}
	xlog.Debug("origin:", string(origin))
}
