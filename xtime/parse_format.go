package xtime

import (
	"strings"
	"time"

	"github.com/iGoogle-ink/gotil"
)

//解析时间
//    时间字符串格式：2006-01-02 15:04:05
func ParseDateTime(timeStr string) (datetime time.Time) {
	datetime, _ = time.ParseInLocation(gotil.TimeLayout, timeStr, time.Local)
	return
}

//解析日期
//    日期字符串格式：2006-01-02
func ParseDate(timeStr string) (date time.Time) {
	date, _ = time.ParseInLocation(gotil.DateLayout, timeStr, time.Local)
	return
}

//格式化Datetime字符串
//    格式化前输入样式：2019-01-04T15:40:00Z 或 2019-01-04T15:40:00+08:00
//    格式化后返回样式：2019-01-04 15:40:00
func FormatDateTime(timeStr string) (formatTime string) {
	if timeStr == gotil.NULL {
		return gotil.NULL
	}
	replace := strings.Replace(timeStr, "T", " ", 1)
	formatTime = replace[:19]
	return
}

//格式化Date成字符串
//    格式化前输入样式：2019-01-04T15:40:00Z 或 2019-01-04T15:40:00+08:00
//    格式化后返回样式：2019-01-04
func FormatDate(dateStr string) (formatDate string) {
	if dateStr == gotil.NULL {
		return gotil.NULL
	}
	split := strings.Split(dateStr, "T")
	formatDate = split[0]
	return
}
