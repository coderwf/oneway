package oneway

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOneWay_Get(t *testing.T) {
    one := NewOneway(NewMapCache())

    var called int64
    c := make(chan string)


    getter := func(key interface{}) (value interface{}, err error) {
    	atomic.AddInt64(&called, 1)
    	return <-c, nil
	}


	const n = 100
	var wait sync.WaitGroup

	for i:= 0; i< n; i++{
		wait.Add(1)
		go func() {
			v, err := one.Get("name", getter)
			if v != "coder"{
				t.Errorf("value: %v", v)
			}
			if err != nil{
				t.Error(err)
			}
			wait.Done()
		}()
	}

    time.Sleep(100 * time.Millisecond)

	c <- "coder"
	wait.Wait()

	if atomic.LoadInt64(&called) != 1{
		t.Errorf("called: %d", called)
	}

    fmt.Println(called)
	fmt.Println(one.Get("name", getter))


}
