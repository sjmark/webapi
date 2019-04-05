package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 获取昨天的起始时间
// 返回结果: 2018-08-29 00:00:00
func GetLastDateTimeBeginString(location *time.Location) string {

	last := time.Now().AddDate(0, 0, -1) // 昨天时间对象

	laststr := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, location).Format("2006-01-02 15:04:05") // 昨天零点字符串

	return laststr
}

// 获取昨天的日期字符串
// 返回结果: 2018-08-29
func GetLastDateString() string {

	location, _ := time.LoadLocation("Asia/Shanghai")

	last := time.Now().AddDate(0, 0, -1) // 昨天时间对象

	laststr := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, location).Format("2006-01-02")

	return laststr
}

func GetLastDateTime() time.Time {

	location, _ := time.LoadLocation("Asia/Shanghai")

	last := time.Now().AddDate(0, 0, -1) // 昨天时间对象

	return time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, location)

}
func GetLastDateTimeWithLocation(location *time.Location) time.Time {

	last := time.Now().AddDate(0, 0, -1) // 昨天时间对象

	return time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, location)

}

// 传入格式:
// "Apr 2, 2018 12:00:00 AM"
// "Mar 29, 2018 12:00:00 AM
// 返回格式:
// 2018-04-02 12:00:00
func TimeParse(timestr string) string {

	if len(timestr) != 24 && len(timestr) != 23 {
		return ""
	}

	var m, d, y, t string

	switch timestr[:3] {
	case "Dec":
		m = "12"
	case "Jan":
		m = "01"
	case "Feb":
		m = "02"
	case "Mar":
		m = "03"
	case "Apr":
		m = "04"
	case "May":
		m = "05"
	case "Jun":
		m = "06"
	case "Jul":
		m = "07"
	case "Aug":
		m = "08"
	case "Sep":
		m = "09"
	case "Oct":
		m = "10"
	case "Nov":
		m = "11"
	}

	if len(timestr) == 24 {
		d = timestr[4:6]
		y = timestr[8:12]
		t = timestr[13:21]
	} else {
		d = "0" + timestr[4:5]
		y = timestr[7:11]
		t = timestr[12:20]
	}

	return fmt.Sprintf("%s-%s-%s %s", y, m, d, t)

}

// 使用说明:
// 输入: 2018-05-12 00:00:00
// 输出: 2018-05-12 23:59:59
func GetLastSecondOfDate(date string, location *time.Location) (string, error) {

	etime, err := time.ParseInLocation("2006-01-02 15:04:05", date, location)

	if err != nil {

		return "", err
	}

	return etime.AddDate(0, 0, 1).Add(-time.Second).Format("2006-01-02 15:04:05"), nil
}

// 日期格式为2006-01-02
func GetDateList(start, end string) ([]string, error) {

	result := make([]string, 0)

	for start <= end {

		result = append(result, start)

		tmp, err := time.Parse("2006-01-02", start)

		if err != nil {
			return nil, err
		}

		start = tmp.AddDate(0, 0, 1).Format("2006-01-02")

	}

	return result, nil

}

// 日期格式为2006-01-02 15:04:05
func GetDateTimeList(start, end string) ([]string, error) {

	result := make([]string, 0)

	for start <= end {

		result = append(result, start)

		tmp, err := time.Parse("2006-01-02 15:04:05", start)

		if err != nil {
			return nil, err
		}

		start = tmp.AddDate(0, 0, 1).Format("2006-01-02 15:04:05")

	}

	return result, nil

}

// 日期格式为2006-01-02 15:00
func FormatTime(revTime string, isHour bool) string {

	var loc, _ = time.LoadLocation("Local")

	var otime, _ = time.ParseInLocation("2006-01-02 15:00", revTime, loc)

	if isHour {
		return otime.Format("2006-01-02 15")
	}

	return otime.Format("2006-01-02")
}

type TimeValue struct {
	Weekly int
	Day    int
	Hour   int
	Minute int
	Second int
}

// 初始化TimeValue的值
func NewTimeValue() *TimeValue {
	return &TimeValue{
		Weekly: 0, // 周
		Day:    0, // 天
		Hour:   0, // 小时
		Minute: 0, // 分
		Second: 0, // 秒
	}
}

/**
*周、天、时、分、秒的处理
*return time.Time
 */
func CommonTime(timeVal *TimeValue) time.Time {
	var next time.Time
	switch {
	case timeVal.Weekly > 0:
		next = time.Now().Add(time.Hour * 24 * 7 * time.Duration(timeVal.Weekly))
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		break
	case timeVal.Day > 0:
		next = time.Now().Add(time.Hour * 24 * time.Duration(timeVal.Hour))
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		break
	case timeVal.Hour > 0:
		next = time.Now().Add(time.Hour * time.Duration(timeVal.Hour))
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
		break
	case timeVal.Minute > 0:
		next = time.Now().Add(time.Minute * time.Duration(timeVal.Minute))
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), 0, 0, next.Location())
		break
	case timeVal.Second > 0:
		next = time.Now().Add(time.Second * time.Duration(timeVal.Second))
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), next.Second(), 0, next.Location())
		break
	default:
		next = time.Now()
	}
	return next
}

func NowHTimeour() int {

	next := time.Now()

	next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())

	s := next.Format("2006-01-02 15")

	times := strings.Split(s, " ")

	if len(times) < 2 {
		return -1
	}

	hour, _ := strconv.Atoi(times[1])

	return hour % 2
}
func GetTimerForMonth(date, hour, minute int) *time.Timer {

	now := time.Now()

	next := now.Add(30 * 24 * time.Hour)
	next = time.Date(next.Year(), next.Month(), date, hour, minute, 0, 0, next.Location())

	return time.NewTimer(next.Sub(now))
}
func GetTimerForWeek(hour, minute int) *time.Timer {

	now := time.Now()

	next := now

	switch now.Weekday() {
	case time.Sunday:
		next = now.Add(1 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	case time.Monday:
		next = now.Add(7 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	case time.Tuesday:
		next = now.Add(6 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	case time.Wednesday:
		next = now.Add(5 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	case time.Thursday:
		next = now.Add(4 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	case time.Friday:
		next = now.Add(3 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	case time.Saturday:
		next = now.Add(2 * 24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	}

	fmt.Println(next.Format("2006-01-02 15:04:05"))

	return time.NewTimer(next.Sub(now))

}
func GetTimerForDate(hour, minute int) *time.Timer {

	now := time.Now()

	var next time.Time

	if now.Hour() < hour || (now.Hour() == hour && now.Minute() < minute) {
		next = now
		next = time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
	} else {
		next = now.Add(24 * time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), hour, minute, 0, 0, next.Location())
	}

	return time.NewTimer(next.Sub(now))

}
func GetTimerForHour(minute int) *time.Timer {

	now := time.Now()

	var next time.Time

	if now.Minute() < minute {
		next = now
		next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), minute, 0, 0, now.Location())
	} else {
		next = now.Add(time.Hour)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), minute, 0, 0, next.Location())
	}

	return time.NewTimer(next.Sub(now))

}
func GetTimerForMinute(second int) *time.Timer {

	now := time.Now()

	var next time.Time

	if now.Second() < second {
		next = now
		next = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), second, 0, now.Location())
	} else {
		next = now.Add(time.Minute)
		next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), next.Minute(), second, 0, next.Location())
	}

	return time.NewTimer(next.Sub(now))

}

//传入格式00:00:15.02获取秒
func TimeForSecond(timestr string) float64 {

	str := strings.Split(timestr, ":")

	h, err := strconv.ParseFloat(str[0], 10)
	if err != nil {
		fmt.Println(err.Error())
	}
	m, err := strconv.ParseFloat(str[1], 10)
	if err != nil {
		fmt.Println(err.Error())
	}
	s, err := strconv.ParseFloat(str[2], 10)
	sec := (h+m)*60 + s

	return sec
}
