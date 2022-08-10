package timeutils

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/juicesix/logging"
	uuid "github.com/satori/go.uuid"
)

const (
	TIME_BASE_FORMAT_DATE = "2006-01-02 15:04:05"
	TIME_BASE_FORMAT_DAY  = "2006-01-02"
)

// MondayTime 获取某个时间的本周周一(不带格式)-返回东八区的时间
func MondayTime(now time.Time) time.Time {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	return GetZeroTime(now.AddDate(0, 0, offset))
}

// MondayTimeLastWeek 获取某个时间的上周周一(不带格式)-返回东八区的时间
func MondayTimeLastWeek(now time.Time) time.Time {
	return MondayTime(now.Add(-86400 * time.Second * 7))
}

// MondayTimeNextWeek 获取某个时间的下周周一(不带格式)-返回东八区的时间
func MondayTimeNextWeek(now time.Time) time.Time {
	return MondayTime(now.Add(86400 * time.Second * 7))
}

// GetZeroTime 获取0点0时0分的时间-返回东八区的时间
func GetZeroTime(now time.Time) time.Time {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println(err)
	}

	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
}

// StrToTime 解析时间字符串(相当于PHP的strtotime，可以指定时区)
func StrToTime(str string) (int64, error) {
	// 使用time.ParseInLocation(使用的时区为：UTC+8)
	// 返回的结果：在UTC+8时区中，等于str这个时间对应的时间戳
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0, err
	}
	timeObj, err := time.ParseInLocation(TIME_BASE_FORMAT_DATE, str, loc)
	if err != nil {
		return 0, err
	}

	return timeObj.Unix(), nil
}

// Date 根据时间戳，返回东八区的时间字符串
func Date(format string, timestamp int64) (string, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return "", err
	}
	return time.Unix(timestamp, 0).In(loc).Format(format), nil
}

func SetDefaultLocation() (*time.Location, error) {
	return time.LoadLocation("Asia/Shanghai")
}

// TimeToStr 某日的时间字符串
func TimeToStr(now time.Time, format ...interface{}) string {
	loc, err := SetDefaultLocation()
	if err != nil {
		return ""
	}
	var f string
	if len(format) == 0 {
		f = TIME_BASE_FORMAT_DAY
	} else {
		f = format[0].(string)
	}
	return now.In(loc).Format(f)
}

func MillisecondToTime(e int64) *time.Time {
	datetime := time.Unix(e/1000, 0)
	return &datetime
}

func TimeToMillisecond(e time.Time) int64 {
	return e.UnixNano() / 1e6
}

func TimeCost(start time.Time, execute, name string) time.Duration {
	terminal := time.Since(start)
	//fmt.Println(fmt.Sprintf("%v-%v 方法耗时:%v", execute, name, terminal))
	logging.Debugf(fmt.Sprintf("%v-%v 方法耗时:%v", execute, name, terminal))
	return terminal
}

//获取当天日期
func NowDateStr() string {
	return time.Now().Format("2006-01-02")
}

func NowTimeStr() string {
	return time.Now().Format("15:04:05")
}

//获取当天日期 不带"-"
func NowDateWithoutLineStr() string {
	return time.Now().Format("20060102")
}

//获取本月 不带"-"
func NowMonthWithoutLineStr() string {
	return time.Now().Format("200601")
}

//一天日期为基准  获取m天后的日期字符串  比如今天20200728,m=1,那么返回20200729，如果m=-1， 返回20200727
func GetDateWithoutLine(m int) string {
	return time.Now().AddDate(0, 0, m).Format("20060102")
}

//获取当前时间戳
func GetNowTimeUnix() int64 {
	return time.Now().Unix()
}

//获取两个数之间的随机数
func RandInt(min int, max int) int {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

//数组取差集
func ExcludeRepeatElem(arr1 []int64, arr2 []int64) []int64 {

	if len(arr1) == 0 {
		return nil
	}

	if len(arr2) == 0 {
		return arr1
	}

	excludeMap := make(map[int64]int64, 0) //利用map去重
	for _, v := range arr1 {
		excludeMap[v] = v
	}

	for _, v := range arr2 {
		if _, ok := excludeMap[v]; ok {
			delete(excludeMap, v)
		}
	}

	ret := make([]int64, 0)
	for _, v := range excludeMap {
		ret = append(ret, v)
	}

	return ret
}

func GetUuid() string {
	uu := uuid.NewV4()
	return uu.String()
}

func GetUuidWithoutMiddleLine() string {
	uu := uuid.NewV4()
	str := strings.ReplaceAll(uu.String(), "-", "")
	return str
}

//"1,2,3,4"  ==> []int{1,2,3,4}
func SplitStringToInts(str string) []int {
	strs := strings.Split(str, ",")
	ret := make([]int, 0)
	for _, v := range strs {
		r, err := strconv.Atoi(v)
		if err != nil {
			logging.Errorf("utils SplitStringToInts fail, v=%s, err=%v", v, err)
			return []int{}
		}
		ret = append(ret, r)
	}

	return ret
}

func Concat(vals []int) string {
	buff := new(bytes.Buffer)
	for _, v := range vals {
		buff.WriteString(fmt.Sprintf(",%d", v))
	}
	return buff.String()[1:]
}

//获取当天零点时间戳
func GetTodayZeroTime() int64 {
	str := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", str, time.Local)
	return t.Unix()
}

//获取当天最后一秒时间戳
func GetTodayLastTime() int64 {
	str := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", str+" 23:59:59", time.Local)
	return t.Unix()
}
