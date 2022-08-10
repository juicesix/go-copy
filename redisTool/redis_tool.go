package redisTool

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/juicesix/logging"
	"github.com/juicesix/redisgogogo"
)

func RedisConn(service Redisgogogo.RedisConfig) *Redisgogogo.Redis {
	max := 100
	for i := 0; i < max; i++ {
		//r, err := rpc.GetRedis(service)
		r, err := Redisgogogo.NewRedis(&service)
		if err != nil {
			if i < max-1 {
				logging.Errorf("redis is nil: %v %#v\n", service, err)
				continue
			}
			logging.Fatalf("redis is nil: %v %#v\n", service, err)
		}
		return r
	}
	return nil
}

type MemberWithScore struct {
	Member string
	Score  int
}

func ZRevRangeWithScore(r *Redisgogogo.Redis, key string, start, end int) ([]MemberWithScore, error) {
	memWithScore, err := r.ZRevRangeWithScore(key, start, end)
	if err != nil {
		return nil, err
	}
	mems := make([]MemberWithScore, 0)
	var member string
	var sc int
	for index, v := range memWithScore {
		if index%2 == 0 {
			member = v
			sc = 0
		} else {
			sc, _ = strconv.Atoi(v)
			mems = append(mems, MemberWithScore{
				Member: member,
				Score:  sc,
			})
		}
	}
	return mems, nil
}

// 简化pipeline流程
func demoPipelining() {
	fn := func(logHead string, pipeline *Redisgogogo.Pipelining, keys []string) {
		for _, key := range keys {
			if err := pipeline.Send("SET", key, 100); err != nil {
				logging.Errorf(logHead+"pipeline.Send, error=%v", err)
			}
		}

		// 把缓冲区中的内容写入到网络
		if err := pipeline.Flush(); err != nil {
			logging.Errorf(logHead+"pipeline.Flush, error=%v", err)
		}

		for _, key := range keys {
			reply, err := redis.String(pipeline.Receive())
			logging.Infof(logHead+"pipeline.Receive: key=%v,reply=%v,err=%v", key, reply, err)
			fmt.Printf(logHead+"pipeline.Receive: key=%v,reply=%v,err=%v\n", key, reply, err)
		}
	}

	// get redis client
	logHead := "logHead|"
	//redisClient, err := rpc.GetRedis("demo")
	redisClient, err := Redisgogogo.NewRedis(&Redisgogogo.RedisConfig{})
	if err != nil {
		logging.Errorf(logHead+"rpc.GetRedis, error=%v", err)
	}
	keys := []string{"a", "b", "c", "d", "e"}
	PipelineMiddleware(context.TODO(), redisClient, logHead, reflect.ValueOf(fn), []reflect.Value{
		reflect.ValueOf(keys),
	})
}

func PipelineMiddleware(ctx context.Context, redisClient *Redisgogogo.Redis, logHead string, fn reflect.Value, args []reflect.Value) (retArr []reflect.Value) {
	// create pipeline
	pipeline, err := redisClient.NewPipelining(ctx)
	if err != nil {
		logging.Errorf(logHead+"redisClient.NewPipelining, error=%v", err)
		return
	}
	defer pipeline.Close()

	// rebuild args
	var allArgs []reflect.Value
	allArgs = append(allArgs, reflect.ValueOf(logHead))
	allArgs = append(allArgs, reflect.ValueOf(pipeline))
	allArgs = append(allArgs, args...)

	return fn.Call(allArgs)
}
