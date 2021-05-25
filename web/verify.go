package web

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iGoogle-ink/gotil/ecode"
	"github.com/iGoogle-ink/gotil/util"
)

const (
	sign  = "sign"
	ts    = "ts"
	appid = "appid"
)

// VerifySign 验证Sign，签名规则，base64(md5(appid+path+ts))
// todo: 当前支持60s验证，后续改成单次验证后即失效
func VerifySign(c *gin.Context) {
	ts := c.GetHeader(ts)
	tsTime := time.Unix(util.String2Int64(ts), 0)
	if time.Now().Sub(tsTime).Seconds() > 60 {
		JSON(c, nil, ecode.InvalidSignErr)
		c.Abort()
		return
	}
	sign := c.GetHeader(sign)
	appid := c.GetHeader(appid)
	path := c.Request.RequestURI
	split := strings.Split(path, "?")
	if !CheckSign(sign, appid, split[0], ts) {
		JSON(c, nil, ecode.InvalidSignErr)
		c.Abort()
		return
	}
}

func CheckSign(sign, appid, path, ts string) (ok bool) {
	hash := md5.New()
	hash.Write([]byte(appid))
	hash.Write([]byte(path))
	hash.Write([]byte(ts))
	calSign := hex.EncodeToString(hash.Sum(nil))
	return sign == calSign
}
