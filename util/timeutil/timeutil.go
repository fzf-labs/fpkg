package timeutil

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/pkg/errors"
)

// TimeLayout 常用日期格式化模板
var TimeLayout = "2006-01-02 15:04:05"

func Carbon() carbon.Carbon {
	return carbon.NewCarbon().SetTimezone(carbon.PRC)
}
func Time2Carbon(t time.Time) carbon.Carbon {
	return carbon.Time2Carbon(t).SetTimezone(carbon.PRC)
}

func NowCarbon() carbon.Carbon {
	return carbon.Now().SetTimezone(carbon.PRC)
}

func NowTime() time.Time {
	return carbon.Now().Carbon2Time()
}

func NowUnix() int64 {
	return carbon.Now().Timestamp()
}

// NowString 转换为当前时间 2021-06-29 23:53:32
func NowString() string {
	return carbon.Now().ToDateTimeString()
}

// NowMillisecondString 转换为当前时间 2021-06-29 23:53:32.010
func NowMillisecondString() string {
	now := carbon.Now()
	return now.ToDateTimeString() + "." + now.Format("u")
}

// NowMicrosecondString 转换为当前时间 2021-06-29 23:53:32.100000
func NowMicrosecondString() string {
	now := carbon.Now()
	return now.ToDateTimeString() + "." + strconv.Itoa(now.Microsecond())
}

// TimeToShortString 时间转日期
func TimeToShortString(ts time.Time) string {
	return time.Unix(ts.Unix(), 00).Format("2006.01.02")
}

// TimeToYearMonthString 时间转日期
func TimeToYearMonthString(ts time.Time) string {
	return time.Unix(ts.Unix(), 00).Format("200601")
}

func ToDateTimeStringByTime(ts time.Time) string {
	return carbon.Time2Carbon(ts).ToDateTimeString()
}

func ToDateTimeStringByTimePointer(ts *time.Time) string {
	if ts != nil {
		return carbon.Time2Carbon(*ts).ToDateTimeString()
	} else {
		return ""
	}
}

// StrToTime 等同于PHP的strtotime函数
// StrToTime("2020-12-19 14:16:22")
func StrToTime(value string) (time.Time, error) {
	var tt time.Time
	if value == "" {
		return tt, errors.New("value is null")
	}
	l, err := time.LoadLocation("Local")
	if err != nil {
		return tt, errors.New("time loadLocation err")
	}
	layouts := []string{
		"20060102",
		"20060102150405",
		"2006-01-02 15:04:05 -0700 MST",
		"2006-01-02 15:04:05 -0700",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05 -0700 MST",
		"2006/01/02 15:04:05 -0700",
		"2006/01/02 15:04:05",
		"2006-01-02 -0700 MST",
		"2006-01-02 -0700",
		"2006-01-02",
		"2006/01/02 -0700 MST",
		"2006/01/02 -0700",
		"2006/01/02",
		"2006-01-02 15:04:05 -0700 -0700",
		"2006/01/02 15:04:05 -0700 -0700",
		"2006-01-02 -0700 -0700",
		"2006/01/02 -0700 -0700",
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC3339Nano,
		time.Kitchen,
		time.Stamp,
		time.StampMilli,
		time.StampMicro,
		time.StampNano,
	}
	for _, layout := range layouts {
		tt, err = time.ParseInLocation(layout, value, l)
		if err == nil {
			return tt, nil
		}
	}
	return tt, errors.New("strToTime err")
}

// StrToLocalTime 字符串转本地时间
func StrToLocalTime(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, errors.New("value is null")
	}
	zoneName, offset := time.Now().Zone()

	zoneValue := offset / 3600 * 100
	if zoneValue > 0 {
		value += fmt.Sprintf(" +%04d", zoneValue)
	} else {
		value += fmt.Sprintf(" -%04d", zoneValue)
	}

	if zoneName != "" {
		value += " " + zoneName
	}
	return StrToTime(value)
}

// GetShowTime 格式化人类友好时间
func GetShowTime(ts time.Time) string {
	duration := time.Now().Unix() - ts.Unix()
	timeStr := ""
	if duration < 60 {
		timeStr = "刚刚发布"
	} else if duration < 3600 {
		timeStr = fmt.Sprintf("%d分钟前更新", duration/60)
	} else if duration < 86400 {
		timeStr = fmt.Sprintf("%d小时前更新", duration/3600)
	} else if duration < 86400*2 {
		timeStr = "昨天更新"
	} else {
		timeStr = TimeToShortString(ts) + "前更新"
	}
	return timeStr
}

// TimeToHuman 根据时间戳获得人类可读时间
func TimeToHuman(ts int) string {
	var res = ""
	if ts == 0 {
		return res
	}

	tt := int(time.Now().Unix()) - ts
	data := []map[string]any{
		{"key": 31536000, "value": "年"},
		{"key": 2592000, "value": "个月"},
		{"key": 604800, "value": "星期"},
		{"key": 86400, "value": "天"},
		{"key": 3600, "value": "小时"},
		{"key": 60, "value": "分钟"},
		{"key": 1, "value": "秒"},
	}
	for _, v := range data {
		var c = tt / v["key"].(int)
		if c != 0 {
			suffix := "前"
			if c < 0 {
				suffix = "后"
				c = -c
			}
			res = strconv.Itoa(c) + v["value"].(string) + suffix
			break
		}
	}

	return res
}

// ToSQLNullTime 将time.Time转换为sql.NullTime
func ToSQLNullTime(tt time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  tt,
		Valid: true,
	}
}

func NowSQLNullTime() sql.NullTime {
	return sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
}
