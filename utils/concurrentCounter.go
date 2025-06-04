package utils

import (
	"sync"
	"sync/atomic"
)

type ConcurrentCounter struct {
	m sync.Map // map[string]*int64
}

func (cc *ConcurrentCounter) Inc(key string) {
	val, _ := cc.m.LoadOrStore(key, new(uint64))
	counter := val.(*uint64)
	atomic.AddUint64(counter, 1)
}

func (cc *ConcurrentCounter) Get(key string) uint64 {
	val, ok := cc.m.Load(key)
	if !ok {
		return 0
	}
	return atomic.LoadUint64(val.(*uint64))
}

func (cc *ConcurrentCounter) Snapshot() map[string]uint64 {
	snapshot := make(map[string]uint64)
	cc.m.Range(func(k, v any) bool {
		snapshot[k.(string)] = atomic.LoadUint64(v.(*uint64))
		return true
	})
	return snapshot
}
