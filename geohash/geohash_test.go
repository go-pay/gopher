package geohash

import (
	"testing"

	"github.com/iGoogle-ink/gopher/xlog"
)

func TestEncode(t *testing.T) {

	geohash := Encode(31.2851847116, 121.5571761131, 10)
	xlog.Debug(geohash) // wtw3yp71rm
}
