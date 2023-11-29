package geohash

import (
	"bytes"
	"strings"

	"github.com/go-pay/util/convert"
)

// geohash精度的设定参考 http://en.wikipedia.org/wiki/Geohash
// geohash length	lat bits	lng bits	lat error	lng error	km error
// 1				2			3			±23			±23			±2500
// 2				5			5			± 2.8		± 5.6		±630
// 3				7			8			± 0.70		± 0.7		±78
// 4				10			10			± 0.087		± 0.18		±20
// 5				12			13			± 0.022		± 0.022		±2.4
// 6				15			15			± 0.0027	± 0.0055	±0.61
// 7				17			18			±0.00068	±0.00068	±0.076
// 8				20			20			±0.000085	±0.00017	±0.019
type Box struct {
	MaxLat float64 //最大纬度
	MinLat float64 //最小纬度
	MaxLng float64 //最大经度
	MinLng float64 //最小经度
}

// 计算宽度
func (b *Box) Width() (width float64) {
	width = b.MaxLng - b.MinLng
	return
}

// 计算高度
func (b *Box) Height() (height float64) {
	height = b.MaxLat - b.MinLat
	return
}

// 计算geohash值
// lat:纬度
// lng:经度
// precision:精度值
func Encode(lat, lng float64, precision int) (geohashCode string) {
	var (
		buffer    = new(bytes.Buffer)
		Base32    = [32]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "b", "c", "d", "e", "f", "g", "h", "j", "k", "m", "n", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
		maxLat    = 90.0              // 初始化最大纬度
		minLat    = -90.0             // 初始化最小纬度
		maxLng    = 180.0             // 初始化最大经度
		minLng    = -180.0            // 初始化最小经度
		mid       float64             // 中间值标记
		bit       = 0                 // 二进制个数标记
		preLen    = 0                 // 初始化计算精度的位数
		isEvenNum = true              // 是否为偶数
		tempBits  = make([]string, 0) // 初始化临时bit数组
	)

	for preLen < precision {
		if isEvenNum {
			//偶数位放经度
			mid = (maxLng + minLng) / 2
			//fmt.Printf("偶数：lng:%v,mid:%v \n", lng, mid)
			if lng < mid {
				maxLng = mid
				tempBits = append(tempBits, "0")
			} else {
				minLng = mid
				tempBits = append(tempBits, "1")
			}
		} else {
			//奇数位放纬度
			mid = (maxLat + minLat) / 2
			//fmt.Printf("奇数：lat:%v,mid:%v \n", lat, mid)
			if lat < mid {
				maxLat = mid
				tempBits = append(tempBits, "0")
			} else {
				minLat = mid
				tempBits = append(tempBits, "1")
			}
		}
		isEvenNum = !isEvenNum
		if bit < 4 {
			bit++
		} else {
			bitNum := strings.Join(tempBits, "")
			num := convert.BinaryToDecimal(bitNum)
			//fmt.Printf("%v => %v => %v\n", tempBits, num, Base32[num])
			buffer.WriteString(Base32[num])
			preLen++
			//重置数据
			bit = 0
			tempBits = make([]string, 0)
		}
	}
	geohashCode = buffer.String()
	return
}
