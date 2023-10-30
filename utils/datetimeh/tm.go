/*
File Name:  tm.py
Description: 默认时区 东8   utc+8   gmt+8
Author:      Chenghu
Date:       2023/9/5 14:29
Change Activity:
*/
package datetimeh

import (
	"github.com/golang-module/carbon/v2"
	"github.com/shopspring/decimal"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

var TimeZonesCarbon = map[int64]carbon.Carbon{}

func timeZoneHandler(zone ...string) carbon.Carbon {
	var offset int64 = 8 * 60 * 60
	if len(zone) > 0 {
		setof, err := strconv.ParseFloat(strings.TrimPrefix(zone[0], "GMT"), 64)
		if err != nil {
			slog.Error("gmt 时间处理失败", "error", err.Error())
		}
		offset = decimal.NewFromFloat(setof * 60 * 60).IntPart()
	}
	if ca, ok := TimeZonesCarbon[offset]; ok {
		return ca
	}
	// 设置时区
	cab := carbon.SetLocation(time.FixedZone("", int(offset)))
	// 设置一周的开始时间为周一
	cab = cab.Now()
	cab = cab.SetWeekStartsAt("Monday")

	TimeZonesCarbon[offset] = cab
	return cab
}

func Now(zone ...string) carbon.Carbon {
	cab := timeZoneHandler(zone...)
	now := cab.Now()
	return now
}

// 今天剩余的时间 秒

func DayNextSecond(zone ...string) int64 {
	// zone 格林时间 GMT+8
	cab := timeZoneHandler(zone...)
	now := cab.Now()
	remainTime := now.DiffInSeconds(now.EndOfDay())
	if remainTime <= 1 {
		return 1
	}
	return remainTime
}

// 今天开始后的时间 秒

func DayLastSecond(zone ...string) int64 {
	// zone 格林时间 GMT+8
	cab := timeZoneHandler(zone...)
	now := cab.Now()
	remainTime := now.StartOfDay().DiffInSeconds(now)
	return remainTime
}

// 本周剩余的时间   周一为开始时间， 周末为结束时间

func WeekNextSecond(zone ...string) int64 {
	cab := timeZoneHandler(zone...)
	now := cab.Now()
	remainTime := now.DiffInSeconds(now.EndOfWeek())
	if remainTime <= 1 {
		return 1
	}
	return remainTime
}

// 获取本月的剩余时间
func MouthNextSexond(zone ...string) int64 {
	cab := timeZoneHandler(zone...)
	remainTime := cab.Now().DiffInSeconds(cab.EndOfMonth())
	if remainTime <= 1 {
		return 1
	}
	return remainTime
}

// 获取当前小时的剩余时间
func HourNextSecond(zone ...string) int64 {
	cab := timeZoneHandler(zone...)
	remainTime := cab.Now().DiffInSeconds(cab.EndOfHour())
	if remainTime <= 1 {
		return 1
	}
	return remainTime
}

// 获取本月的剩余时间
func ConvertTimestampToDateTime(timestamp int64, zone ...string) carbon.Carbon {
	cab := timeZoneHandler(zone...)
	cab = cab.CreateFromTimestamp(timestamp, zone...)
	return cab
}
