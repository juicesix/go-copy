package lru_cache

import (
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/juicesix/logging"
)

func GetLruCacheKeyCv(uid int64) string {
	return fmt.Sprintf("lru_cache_cv_%d", uid)
}

const (
	LruCacheSize = 50000
)

var (
	LruCacheOnce sync.Once
	LruCache     *lru.Cache
)

func InitLruCache() {
	LruCacheOnce.Do(func() {
		c, err := lru.New(LruCacheSize)
		if err != nil {
			panic(fmt.Sprintf("InitProfileCache failed error %v", err))
		}
		LruCache = c
	})
}

var (
	cacheHitNum float64
	visitNum    float64
)

// UpdateUserCvLruCache 更新LRU缓存
func UpdateUserCvLruCache(uids []int64, module string) {
	logHead := fmt.Sprintf("UpdateUserCvLruCache(%s)|", module)
	var getFromApi []int64

	visitNum += float64(len(uids))

	for _, uid := range uids {
		cv, ok := LruCache.Get(GetLruCacheKeyCv(uid))
		logging.Debugf(logHead+"cv=%v,ok=%v", cv, ok)
		if !ok {
			getFromApi = append(getFromApi, uid)
		} else {
			cacheHitNum++
		}
	}
	logging.Debugf(logHead+"getFromApi=%v", getFromApi)
	var hitRatio float64
	if visitNum != 0 {
		hitRatio = cacheHitNum / visitNum
	} else {
		hitRatio = 0
	}
	// 统计缓存命中率
	logging.Infof(logHead+"cacheHitNum=%v,visitNum=%v,hitRatio=%v", cacheHitNum, visitNum, hitRatio)

	// get from api
	if len(getFromApi) > 0 {
		//retUserInfos := GetUserInfoByUids(context.Background(), getFromApi)
		//for _, item := range retUserInfos {
		//	if item.ID > 0 {
		//		logging.Debugf(logHead+"Add cache uid=%v,cv=%v", item.ID, item.CV)
		//		LruCache.Add(GetLruCacheKeyCv(item.ID), item.CV)
		//	}
		//}
	}
}
