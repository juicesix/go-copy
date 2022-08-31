package array

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/spf13/cast"
)

var TrimZero = regexp.MustCompile("0+?$")

var TrimPoint = regexp.MustCompile("[.]$")

func Array2StringInt64(arr []int64) string {
	var res string
	for _, item := range arr {
		if "" == res {
			res += fmt.Sprintf("%d", item)
		} else {
			res += fmt.Sprintf(",%d", item)
		}
	}
	return res
}

func Array2String(arr []string) string {
	var res string
	for _, item := range arr {
		if "" == res {
			res += fmt.Sprintf("%s", item)
		} else {
			res += fmt.Sprintf(",%s", item)
		}
	}
	return res
}

// InArray in_array()
// haystack supported types: slice, array or map
// 注意，needle和haystack的类型一定要保证相同，否则会有问题！！！
func InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type muset be slice, array or map")
	}

	return false
}

func ArraySlice(s []int64, offset, length int) []int64 {
	if offset > len(s) {
		return []int64{}
	}
	end := offset + length
	if end < len(s) {
		return s[offset:end]
	}
	return s[offset:]
}

// InSliceIface checks given interface in interface slice.
func InSliceInt64(v int64, sl []int64) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}

func HasChinese(s string) bool {
	for _, v := range s {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

func GetStrEnd(s string) string {
	lastStr := s[len(s)-1:]
	return lastStr
}

func GetUidEnd(uid int64) int64 {
	if uid <= 0 {
		return 0
	}
	lastNum := GetStrEnd(cast.ToString(uid))
	return cast.ToInt64(lastNum)
}

// RmDecimalZero 去除小数位无用的 0
func RmDecimalZero(s string) string {
	if strings.Index(s, ".") <= 0 {
		return s
	}
	// 去掉后面无用的零
	s = strings.ReplaceAll(s, TrimZero.FindString(s), "")
	// 如小数点后面全是零则去掉小数点,
	s = strings.ReplaceAll(s, TrimPoint.FindString(s), "")
	return s
}

func AddArray(target, source []int64) []int64 {
	for _, vv := range source {
		target = append(target, vv)
	}
	return target
}
