package engine_test

import (
	"sync"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	wg := sync.WaitGroup{}

	for range "..." {
		wg.Add(1)
		go func() {
			defer wg.Done()
			t.Log("=====> go func")
			time.Sleep(time.Second * 5)
			t.Log("go func done! <<<<========")
		}()
	}

	waitDone := func() <-chan struct{} {
		wg.Wait()
		ch := make(chan struct{})
		defer close(ch)
		return ch
	}

	select {
	case <-waitDone():
		t.Log("all goroutine done")
	case <-time.After(time.Second * 10):
		t.Error("timeout")
	}
}
