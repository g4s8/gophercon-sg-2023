package main

import (
	"testing"
)

func BenchmarkInt32(b *testing.B) {
	b.Run("1/toInt", newInt32ToIntBenchmark(1))
	b.Run("1/toInt32", newInt32ToInt32Benchmark(1))
	b.Run("1/toIntGeneric", newInt32ToIntGenericBenchmark(1))

	b.Run("256/toInt", newInt32ToIntBenchmark(256))
	b.Run("256/toInt32", newInt32ToInt32Benchmark(256))
	b.Run("256/toIntGeneric", newInt32ToIntGenericBenchmark(256))
}

func BenchmarkInt64(b *testing.B) {
	b.Run("1/toInt", newInt64ToIntBenchmark(1))
	b.Run("1/toInt64", newInt64ToInt64Benchmark(1))
	b.Run("1/toIntGeneric", newInt64ToIntGenericBenchmark(1))

	b.Run("256/toInt", newInt64ToIntBenchmark(256))
	b.Run("256/toInt64", newInt64ToInt64Benchmark(256))
	b.Run("256/toIntGeneric", newInt64ToIntGenericBenchmark(256))
}

func BenchmarkStr(b *testing.B) {
	b.Run("toInt", newStrToIntBenchmark("1"))
	b.Run("toIntStr", newStrToIntStrBenchmark("1"))
	b.Run("toIntGeneric", newStrToIntGenericBenchmark("1"))
}

func BenchmarkStruct(b *testing.B) {
	b.Run("toInt", newStructToIntBenchmark(1))
	b.Run("toIntStruct", newStructToIntStructBenchmark(1))
	b.Run("toIntGeneric", newStructToIntGenericBenchmark(1))
}

func newInt32ToIntBenchmark(val int32) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := inter32(val)
			toInt(target)
		}
	}
}

func newInt32ToInt32Benchmark(val int32) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := inter32(val)
			toInt32(target)
		}
	}
}

func newInt32ToIntGenericBenchmark(val int32) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := inter32(val)
			toIntGeneric(target)
		}
	}
}

func newInt64ToIntBenchmark(val int64) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := inter64(val)
			toInt(target)
		}
	}
}

func newInt64ToInt64Benchmark(val int64) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := inter64(val)
			toInt64(target)
		}
	}
}

func newInt64ToIntGenericBenchmark(val int64) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := inter64(val)
			toIntGeneric(target)
		}
	}
}

func newStrToIntBenchmark(val string) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := interStr(val)
			toInt(target)
		}
	}
}

func newStrToIntStrBenchmark(val string) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := interStr(val)
			toIntStr(target)
		}
	}
}

func newStrToIntGenericBenchmark(val string) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := interStr(val)
			toIntGeneric(target)
		}
	}
}

func newStructToIntBenchmark(val int64) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := interStruct{x: val}
			toInt(&target)
		}
	}
}

func newStructToIntStructBenchmark(val int64) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := interStruct{x: val}
			toIntStruct(&target)
		}
	}
}

func newStructToIntGenericBenchmark(val int64) func(*testing.B) {
	return func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			target := interStruct{x: val}
			toIntGeneric(&target)
		}
	}
}
