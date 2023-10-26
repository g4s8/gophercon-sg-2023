package pools

import (
	"sync"
	"testing"
)

const K8 = 8 * 1024

func BenchmarkChanPool(b *testing.B) {
	p := NewChanPool(K8, 1024)
	b.Run("1", benchmarkPool(p, 1))
	b.Run("2", benchmarkPool(p, 2))
	b.Run("4", benchmarkPool(p, 4))
	b.Run("8", benchmarkPool(p, 8))
	b.Run("16", benchmarkPool(p, 16))
	b.Run("32", benchmarkPool(p, 32))
	b.Run("64", benchmarkPool(p, 64))
	b.Run("128", benchmarkPool(p, 128))
	b.Run("256", benchmarkPool(p, 256))
	b.Run("512", benchmarkPool(p, 512))
	b.Run("1024", benchmarkPool(p, 1024))
}

func BenchmarkSyncPool(b *testing.B) {
	p := NewSyncPool(K8)
	b.Run("1", benchmarkPool(p, 1))
	b.Run("2", benchmarkPool(p, 2))
	b.Run("4", benchmarkPool(p, 4))
	b.Run("8", benchmarkPool(p, 8))
	b.Run("16", benchmarkPool(p, 16))
	b.Run("32", benchmarkPool(p, 32))
	b.Run("64", benchmarkPool(p, 64))
	b.Run("128", benchmarkPool(p, 128))
	b.Run("256", benchmarkPool(p, 256))
	b.Run("512", benchmarkPool(p, 512))
	b.Run("1024", benchmarkPool(p, 1024))
}

type pool interface {
	Get() buffer
	Put(buffer)
	Stats() *Stats
}

func benchmarkPool(p pool, workers int) func(*testing.B) {
	return func(t *testing.B) {
		var wg sync.WaitGroup
		wg.Add(workers)
		p.Stats().Reset()
		for w := 0; w < workers; w++ {
			go func() {
				defer wg.Done()

				for n := 0; n < t.N/workers; n++ {
					buf := p.Get()
					fill(buf, n+w*100)
					workload(buf)
					p.Put(buf)
				}
			}()
		}
		wg.Wait()
		// t.Logf("Stats: %+v", p.Stats())
	}
}

func BenchmarkBaseline(b *testing.B) {
	buf := make(buffer, K8)
	for n := 0; n < b.N; n++ {
		fill(buf, n+1)
		workload(buf)
	}
}

func fill(b buffer, n int) {
	for i := range b {
		b[i] = byte(i % n)
	}
}

const workloadIterations = 100

func workload(b buffer) uint8 {
	var hash uint8
	for i := 0; i < workloadIterations; i++ {
		for i := range b {
			hash += 13 * uint8(i^int(b[i]))
		}
	}
	return hash
}
