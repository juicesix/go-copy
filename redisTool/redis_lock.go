package redisTool

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/juicesix/logging"
	Redisgogogo "github.com/juicesix/redisgogogo"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cast"
)

var uuidClient uuid.UUID

var REDIS_LOCK_EXPIRE = 30

var REDIS_LOCK_EXPIRE_MIN = 3

const (
	REDIS_SUCCESS         = -100
	REDIS_VALUE_EXIST     = -99
	REDIS_VALUE_NOT_EXIST = -98
	REDIS_ILLEGAL_PARAM   = -97
)

const (
	SCRIPT_LOCK = ` 
    local res=redis.call('GET', KEYS[1])
    if res then
        return -99
    else
        redis.call('SET',KEYS[1],ARGV[1]);
        redis.call('EXPIRE',KEYS[1],ARGV[2])
        return -100
    end 
    `

	SCRIPT_EXPIRE = ` 
    local res=redis.call('GET', KEYS[1])
    if not res then
        return -98
    end 
    if res==ARGV[1] then
        redis.call('EXPIRE', KEYS[1], ARGV[2])
        return -100
    else
        return -97
    end 
    `

	SCRIPT_DEL = ` 
    local res=redis.call('GET', KEYS[1])
    if not res then 
        return -98
    end 
    if res==ARGV[1] then
        redis.call('DEL', KEYS[1])
        return -100
    else
        return -97
    end 
    `
)

type RedisLockParam struct {
	Mu      sync.Mutex
	Redis   *Redisgogogo.Redis
	Key     string
	UuidStr string
	Expire  int
}

func Uuid() string {
	u := uuid.NewV4()
	uuidStr := strings.Replace(u.String(), "-", "", -1)
	return uuidStr
}

func LockRedis(r *Redisgogogo.Redis, key string) (bool, *RedisLockParam) {
	param := &RedisLockParam{
		Mu:    sync.Mutex{},
		Redis: r,
		Key:   key,
	}
	param.UuidStr = Uuid()
	isOk := false
	result, err := redis.Int(param.Redis.Eval(SCRIPT_LOCK, []string{param.Key}, []interface{}{param.UuidStr, REDIS_LOCK_EXPIRE}))
	if err != nil {
		logging.Errorf("redis_lock Lock error, key:[%v] , err:[%v] ", param.Key, err)
	}
	if result == REDIS_SUCCESS {
		param.Expire = REDIS_LOCK_EXPIRE
		//go Renew(param)
		isOk = true
	}
	if result == REDIS_VALUE_EXIST {
		logging.Errorf("redis_lock lock failed,key:[%v]", param.Key)
	}
	return isOk, param
}

func UnLockRedis(param *RedisLockParam) bool {
	isOk := false
	result, err := redis.Int(param.Redis.Eval(SCRIPT_DEL, []string{param.Key}, []interface{}{param.UuidStr}))
	if err != nil {
		logging.Errorf("redis_lock UnLock error, key:[%v] , err:[%v],uuidStr:[%v] ", param.Key, err, param.UuidStr)
		return false
	}
	if result != REDIS_SUCCESS {
		switch result {
		case REDIS_VALUE_NOT_EXIST:
			logging.Errorf("redis_lock UnLockRedis failed, key:[%v] , result:[%v],uuidStr:[%v] ", param.Key, REDIS_VALUE_NOT_EXIST, param.UuidStr)
		case REDIS_ILLEGAL_PARAM:
			logging.Errorf("redis_lock UnLockRedis failed, key:[%v] , result:[%v] ,uuidStr:[%v] ", param.Key, REDIS_ILLEGAL_PARAM, param.UuidStr)
		default:
			logging.Errorf("redis_lock UnLockRedis failed, key:[%v] , result:[%v] ,uuidStr:[%v] ", param.Key, result, param.UuidStr)
		}
	} else {
		logging.Infof("redis_lock UnLockRedis success, key:[%v] , result:[%v] ,uuidStr:[%v]", param.Key, result, param.UuidStr)
	}
	return isOk
}

func Renew(param *RedisLockParam) {
	count := 0
	logging.Infof("redis_lock Renew ing, key:[%v],uuidStr:[%v] ,", param.Key, param.UuidStr)
	for true {
		count += 1
		result, err := redis.Int(param.Redis.Eval(SCRIPT_EXPIRE, []string{param.Key}, []interface{}{param.UuidStr, REDIS_LOCK_EXPIRE}))
		if err != nil {
			logging.Errorf("redis_lock Renew error, key:[%v] , err:[%v] ", param.Key, err)
			break
		}
		if result != REDIS_SUCCESS {
			switch result {
			case REDIS_VALUE_NOT_EXIST:
				logging.Errorf("redis_lock Renew failed, key:[%v] , result:[%v],uuidStr:[%v] ", param.Key, REDIS_VALUE_NOT_EXIST, param.UuidStr)
			case REDIS_ILLEGAL_PARAM:
				logging.Errorf("redis_lock Renew failed, key:[%v] , result:[%v] ,uuidStr:[%v] ", param.Key, REDIS_ILLEGAL_PARAM, param.UuidStr)
			default:
				logging.Errorf("redis_lock Renew failed, key:[%v] , result:[%v] ,uuidStr:[%v] ", param.Key, result, param.UuidStr)
			}
			break
		} else {
			logging.Infof("redis_lock Renew success, key:[%v] , result:[%v] ,uuidStr:[%v]", param.Key, result, param.UuidStr)
		}
		if count >= 3 {
			logging.Infof("redis_lock Renew count >= 3 , key:[%v] , result:[%v] ,uuidStr:[%v]", param.Key, result, param.UuidStr)
			break
		}
		time.Sleep(time.Second * 5)
	}
	logging.Infof("redis_lock Renew end, key:[%v] ,uuidStr:[%v] ,", param.Key, param.UuidStr)
}

/**
 * @Author xiaohuihui
 * @Description redis key 设置过期时间
 * @Date 10:19 2021/7/5
 * @Param expireReset:是否直接重置过期时间
 * @Param expireAdd:是否将 expire 累加到过期时间
 * @return
 **/
func RenewSaveExpire(rds *Redisgogogo.Redis, key string, expire time.Duration, expireReset bool, expireAdd bool) error {
	// 直接重置过期时间
	if expireReset {
		rds.Expire(key, expire)
		return nil
	}
	exp, err := redis.Int(rds.Do("ttl", key))
	if err != nil {
		logging.Errorf("RenewSaveExpireV2 redis.Int error:%v,key:[%v]", err, key)
		return err
	}
	// key 不存在时
	if exp == -2 {
		logging.Errorf("RenewSaveExpire exp == -2 , key:%v", key)
		return nil
	}
	expireTime := cast.ToInt64(expire.Seconds())
	expTime := cast.ToInt64(exp)
	// 是否累加过期时间
	if expireAdd {
		err = rds.Expire(key, time.Second*time.Duration(expireTime+expTime))
	}
	// 如果没有过期时间
	if expTime <= -1 {
		err = rds.Expire(key, time.Second*time.Duration(expireTime))
	}
	if err != nil {
		logging.Errorf("RenewSaveExpire rds.Expire error:%v", err)
	}
	return nil
}

func RenewSaveExpireV2(ctx context.Context, rds *Redisgogogo.Redis, key string, expire time.Duration, expireReset bool, expireAdd bool) error {
	// 直接重置过期时间
	if expireReset {
		rds.Expire(key, expire)
		return nil
	}
	exp, err := redis.Int(rds.Do("ttl", key))
	if err != nil {
		logging.Errorf("RenewSaveExpireV2 redis.Int error:%v,key:[%v]", err, key)
		return err
	}
	// key 不存在时
	if exp == -2 {
		logging.Errorf("RenewSaveExpire exp == -2 , key:%v", key)
		return nil
	}
	expireTime := cast.ToInt64(expire.Seconds())
	expTime := cast.ToInt64(exp)
	// 是否累加过期时间
	if expireAdd {
		err = rds.Expire(key, time.Second*time.Duration(expireTime+expTime))
	}
	// 如果当前没有过期时间
	if expTime <= -1 {
		err = rds.Expire(key, time.Second*time.Duration(expireTime))
	}
	if err != nil {
		logging.Errorf("RenewSaveExpireV2 rds.Expire error:%v", err)
	}
	return nil
}

func QueryLock(rds *Redisgogogo.Redis, key string) bool {
	exist, err := rds.Exists(key)
	if err != nil {
		logging.Errorf("QueryLock rds.Exists ,key:%v,error:%v", key, err)
	}
	return exist
}

func LockQuery(rds *Redisgogogo.Redis, key string, limitTime int) (string, error) {
	return rds.SetExSecond(key, key, limitTime)
}
