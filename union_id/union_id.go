package union_id

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cast"
)

type UniId struct {
}

// GenerateRandInt64 生成[min,max)范围内随机数
func GenerateRandInt64(min int64, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(max-min) + min
	return randNum
}

// GenerateRandInt 生成[min,max)范围内随机数
func GenerateRandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(max-min) + min
	return randNum
}

// GenerateRandNum64 生成[min,max)范围内随机数
func (uniId UniId) GenerateRandNum64(min int64, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(max-min) + min
	return randNum
}

// GenerateIntId 适用: 生成uid,具有一定的递增
// 返回值14位，对应数据库bigint类型
func (uniId UniId) GenerateIntId() int64 {
	str := fmt.Sprintf("%v%v", time.Now().Unix(), uniId.GenerateRandNum64(1000, 10000))
	return cast.ToInt64(str)
}

// GenerateIntId64 适用: 生成uid,具有一定的递增
// 返回值19位，对应数据库bigint类型
func (uniId UniId) GenerateIntId64() int64 {
	str := fmt.Sprintf("%v%v", time.Now().UnixNano()/1000, uniId.GenerateRandNum64(100, 1000))
	return cast.ToInt64(str)
}
