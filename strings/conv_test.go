package strings

import (
	"fmt"
	"sort"
	"testing"
)

var casters = map[string]func([]byte) string{
	"string": bytesToString,
	"unsafe": bytesToStringUnsafe,
}

func BenchmarkStringsCast(b *testing.B) {
	b.Skip("skip")
	names := make([]string, 0, len(casters))
	for name := range casters {
		names = append(names, name)
	}
	sort.Strings(names)
	for strLen := 1; strLen <= 1000; strLen *= 10 {
		for rsName, rs := range allSources {
			for _, name := range names {
				f := casters[name]
				b.Run(fmt.Sprintf("%d/%s/%s", strLen, rsName, name),
					benchCastWithParams(strLen, rs, f))
			}
		}
	}
}

func benchCastWithParams(strLen int, rs randomSource, f func([]byte) string) func(*testing.B) {
	return func(b *testing.B) {
		target := []byte(rs.get(0, strLen))
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := f(target)
			if debug {
				b.Logf("res = %q\n", res)
			}
		}
	}
}

func TestCasters(t *testing.T) {
	for strlen := 0; strlen < 100; strlen++ {
		for rsName, rs := range allSources {
			for name, caster := range casters {
				src := []byte(rs.get(0, strlen))
				t.Run(fmt.Sprintf("%d/%s/%s", strlen, rsName, name), func(t *testing.T) {
					dst := caster(src)
					if string(src) != dst {
						t.Errorf("cast failed: %q != %q", src, dst)
					}
				})
			}
		}
	}
}
