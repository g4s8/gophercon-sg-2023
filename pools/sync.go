package pools

import (
	"sync"
	"sync/atomic"
)

type SyncPool struct {
	pool *sync.Pool

	stats Stats
}

func NewSyncPool(bufSize int) *SyncPool {
	var p SyncPool
	p.pool = &sync.Pool{
		New: func() interface{} {
			atomic.AddInt32(&p.stats.Allocated, 1)
			return make(buffer, bufSize)
		},
	}
	return &p
}

func (p *SyncPool) Get() buffer {
	atomic.AddInt32(&p.stats.Used, 1)
	return p.pool.Get().(buffer)
}

func (p *SyncPool) Put(buf buffer) {
	atomic.AddInt32(&p.stats.Recycled, 1)
	buf.reset()
	p.pool.Put(buf)
}

func (p *SyncPool) Stats() *Stats {
	return &p.stats
}
