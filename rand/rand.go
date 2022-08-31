package rand

import (
	"errors"
	"math/rand"
	"time"
)

func Random(lis []interface{}, length int) ([]interface{}, error) {
	rand.Seed(time.Now().Unix())
	if len(lis) <= 0 {
		return []interface{}{}, errors.New("the length of the parameter lis should not be less than 0")
	}

	if length <= 0 { // || len(lis) <= length
		return []interface{}{}, errors.New("the size of the parameter length illegal")
	} else if len(lis) <= length {
		length = len(lis)
	}

	for i := len(lis) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		lis[i], lis[num] = lis[num], lis[i]
	}

	return lis[:length], nil
}
