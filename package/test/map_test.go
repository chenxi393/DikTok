package test

import (
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	userMap := make(map[int]*string)
	for i := range 1000000 {
		s := ""
		userMap[i] = &s
	}
	wg := sync.WaitGroup{}
	wg.Add(len(userMap))
	for i := range userMap {
		go func(id int) {
			defer wg.Done()
			ss := "222"
			*userMap[id] = ss
		}(i)
	}
	wg.Wait()
	// 这样测试好像是不会有并发问题的？
}
