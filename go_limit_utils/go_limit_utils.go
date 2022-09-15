package go_limit_utils

import (
	"context"
	"fmt"
	"sync"
)

// GoLimit 携程限制
type GoLimit struct {
	MD    *sync.WaitGroup
	Max   int
	Do    func(data ...interface{})
	index int
}

func (goLimit *GoLimit) Doing(data ...interface{}) {
	if goLimit.Do == nil {
		return
	}
	if goLimit.Max == 0 {
		go goLimit.Do(data...)
		return
	}
	if goLimit.MD == nil {
		goLimit.MD = &sync.WaitGroup{}
	}
	if goLimit.index == goLimit.Max {
		goLimit.MD.Wait()
		goLimit.index = 0
	}
	if goLimit.index == 0 {
		goLimit.MD.Add(goLimit.Max)
	}
	go func() {
		defer goLimit.MD.Done()
		goLimit.Do(data...)
	}()

	goLimit.index++
}

func GoodBoy(ctx context.Context) {
	go func() {
		limit := GoLimit{
			MD:  &sync.WaitGroup{},
			Max: 10,
			Do: func(data ...interface{}) {
				CallMeGoodBoy(data[0].(context.Context), data[1].(string))
			},
		}
		for i := 0; i < 20; i++ {
			limit.Doing(ctx, "xiao bao")
		}
	}()
}

func CallMeGoodBoy(ctx context.Context, msg string) {
	fmt.Println(msg)
}
