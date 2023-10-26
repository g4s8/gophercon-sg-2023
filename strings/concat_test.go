package strings

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

var joiners = map[string]func([]string) string{
	// "Concat":           concatStringsAdd,
	"Join": concatStringsJoin,
	// "Builder":          concatStringsBuilder,
	"BuilderGrow": concatStringsBuilderGrow,
	// "Buffer":           concatStringsBuffer,
	// "BufferGrow":       concatStringsBufferGrow,
	"Copy":            concatStringsCopy,
	"Append":          concatStringsAppend,
	"CopyBytesUnsafe": concatStringsCopyUnsafe,
	"AppendUnsafe":    concatStringsAppendUnsafe,
	// "CopyBytesUnsafe2": concatStringsCopyUnsafe2,
}

type cachedJoiner interface {
	concat([]string) string
}

var cachedJoiners = map[string]func() cachedJoiner{
	"Builder": func() cachedJoiner { return &builderConcater{} },
	// "Buffer":      func() cachedJoiner { return &bufferConcater{} },
	"Copy":         func() cachedJoiner { return &copyConcater{} },
	"CopyUnsafe":   func() cachedJoiner { return &copyUnsafeConcater{} },
	"AppendUnsafe": func() cachedJoiner { return &appendUnsafeConcater{} },
	// "CopyUnsafe2": func() cachedJoiner { return &copyUnsafe2Concater{} },
}

func BenchmarkStringsConcats(b *testing.B) {
	names := make([]string, 0, len(joiners))
	for name := range joiners {
		names = append(names, name)
	}
	sort.Strings(names)
	for strCount := 10; strCount <= 1000; strCount *= 10 {
		for strLen := 10; strLen <= 1000; strLen *= 10 {
			for rsName, rs := range asciiSources {
				for _, name := range names {
					f := joiners[name]
					b.Run(fmt.Sprintf("%d/%d/%s/%s", strCount, strLen, rsName, name),
						benchConcatWithParams(strCount, strLen, rs, f))
				}
			}
		}
	}
}

func BenchmarkCachedJoiners(b *testing.B) {
	b.Skip("disabled")
	names := make([]string, 0, len(cachedJoiners))
	for name := range cachedJoiners {
		names = append(names, name)
	}
	sort.Strings(names)

	for strCount := 10; strCount <= 1000; strCount *= 10 {
		for strLen := 10; strLen <= 1000; strLen *= 10 {
			for rsName, rs := range asciiSources {
				for _, name := range names {
					f := cachedJoiners[name]
					b.Run(fmt.Sprintf("%d/%d/%s/%s", strCount, strLen, rsName, name),
						benchCachedJoiner(strCount, strLen, rs, f))
				}
			}
		}
	}
}

func benchCachedJoiner(numStrings int, strlen int, rs randomSource, f func() cachedJoiner) func(*testing.B) {
	return func(b *testing.B) {
		ss := make([]string, numStrings)
		for i := 0; i < numStrings; i++ {
			ss[i] = rs.get(i, strlen)
		}

		j := f()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := j.concat(ss)
			if debug {
				b.Logf("res = %q\n", res)
			}
		}
	}
}

func benchConcatWithParams(numStrings int, strLen int, rs randomSource, f func([]string) string) func(*testing.B) {
	return func(b *testing.B) {
		ss := make([]string, numStrings)
		for i := 0; i < numStrings; i++ {
			ss[i] = rs.get(i, strLen)
			if debug {
				b.Logf("ss[%d] = %q\n", i, ss[i])
			}
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			res := f(ss)
			if debug {
				b.Logf("res = %q\n", res)
			}
		}
	}
}

func TestJoiners(t *testing.T) {
	for strCount := 1; strCount <= 1000; strCount *= 10 {
		for strLen := 1; strLen <= 1000; strLen *= 10 {
			for rsName, rs := range allSources {
				ss := make([]string, strCount)
				for i := 0; i < strCount; i++ {
					ss[i] = rs.get(i, strLen)
				}
				for name, f := range joiners {
					t.Run(fmt.Sprintf("%d/%d/%s/%s", strCount, strLen, rsName, name),
						func(t *testing.T) {
							res := f(ss)
							expect := strings.Join(ss, "")
							if res != expect {
								t.Errorf("got %q, expected %q", res, expect)
							}
						})
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestCachedJoiner(t *testing.T) {
	const iters = 100
	for strCount := 1; strCount <= 1000; strCount *= 10 {
		for strLen := 1; strLen <= 1000; strLen *= 10 {
			for rsName, rs := range allSources {
				ss := make([][]string, iters)
				for i := 0; i < iters; i++ {
					ss[i] = make([]string, strCount)
					for j := 0; j < strCount; j++ {
						ss[i][j] = rs.get(i+j, strLen+abs(i-j))
					}
				}
				for name, f := range cachedJoiners {
					t.Run(fmt.Sprintf("%d/%d/%s/%s", strCount, strLen, rsName, name),
						func(t *testing.T) {
							j := f()
							for i := 0; i < iters; i++ {
								res := j.concat(ss[i])
								expect := strings.Join(ss[i], "")
								if res != expect {
									t.Errorf("got %q, expected %q", res, expect)
								}
							}
						})
				}
			}
		}
	}
}
