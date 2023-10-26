package pools

import (
	"strconv"
	"strings"
	"sync/atomic"
)

type buffer []byte

func (b buffer) reset() {
	for i := range b {
		b[i] = 0
	}
}

type Stats struct {
	Allocated int32
	Used      int32
	Reused    int32
	Recycled  int32
	Skipped   int32
}

func (s Stats) String() string {
	var sb strings.Builder
	sb.WriteString("Allocated: ")
	sb.WriteString(strconv.Itoa(int(s.Allocated)))
	sb.WriteString(", Used: ")
	sb.WriteString(strconv.Itoa(int(s.Used)))
	sb.WriteString(", Recycled: ")
	sb.WriteString(strconv.Itoa(int(s.Recycled)))
	if s.Reused > 0 {
		sb.WriteString(", Reused: ")
		sb.WriteString(strconv.Itoa(int(s.Reused)))
	}
	if s.Skipped > 0 {
		sb.WriteString(", Skipped: ")
		sb.WriteString(strconv.Itoa(int(s.Skipped)))
	}
	return sb.String()
}

func (s *Stats) Reset() {
	var empty Stats
	*s = empty
}

type ChanPool struct {
	ch      chan buffer
	bufSize int

	stats Stats
}

func NewChanPool(bufSize int, poolSize int) *ChanPool {
	return &ChanPool{
		ch:      make(chan buffer, poolSize),
		bufSize: bufSize,
	}
}

func (p *ChanPool) Get() buffer {
	atomic.AddInt32(&p.stats.Used, 1)
	select {
	case buf := <-p.ch:
		buf.reset()
		atomic.AddInt32(&p.stats.Reused, 1)
		return buf
	default:
		atomic.AddInt32(&p.stats.Allocated, 1)
		return make(buffer, p.bufSize)
	}
}

func (p *ChanPool) Put(buf buffer) {
	select {
	case p.ch <- buf:
		atomic.AddInt32(&p.stats.Recycled, 1)
	default:
		atomic.AddInt32(&p.stats.Skipped, 1)
	}
}

func (p *ChanPool) WarmUp() {
	for {
		select {
		case p.ch <- make(buffer, p.bufSize):
		default:
			return
		}
	}
}

func (p *ChanPool) Stats() *Stats {
	return &p.stats
}
