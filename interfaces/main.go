package main

type Inter interface {
	Int64() int64
}

func main() {
	// i32_1 := inter32(1)
	// _ = toInt(i32_1)
	// _ = toInt32(i32_1)
	// _ = toIntGeneric(i32_1)

	// i32_256 := inter32(256)
	// _ = toInt(i32_256)
	// _ = toInt32(i32_256)
	// _ = toIntGeneric(i32_256)

	// i64_1 := inter64(1)
	// _ = toInt(i64_1)
	// _ = toInt64(i64_1)
	// _ = toIntGeneric(i64_1)

	// i64_256 := inter64(256)
	// _ = toInt(i64_256)
	// _ = toInt64(i64_256)
	// _ = toIntGeneric(i64_256)

	// iStr := interStr("1")
	// _ = toInt(iStr)
	// _ = toIntStr(iStr)
	// _ = toIntGeneric(iStr)

	// iPtr := interPtr(1)
	// _ = toInt(&iPtr)
	// _ = toIntPtr(&iPtr)
	// _ = toIntGeneric(&iPtr)

	// iStruct := interStruct{1}
	// _ = toInt(iStruct)
	// _ = toIntStruct(iStruct)
	// _ = toIntGeneric(iStruct)

	// iPtrStruct := interPtrStruct{1}
	// _ = toInt(&iPtrStruct)
	// _ = toIntPtrStruct(&iPtrStruct)
	// _ = toIntGeneric(&iPtrStruct)

	// iStructPtr := interStructPtr{}
	// iStructPtr.val = new(int64)
	// *iStructPtr.val = 1
	// _ = toInt(iStructPtr)
	// _ = toIntStructPtr(iStructPtr)
	// _ = toIntGeneric(iStructPtr)

	// iPtrStructPtr := interPtrStructPtr{}
	// iPtrStructPtr.val = new(int64)
	// *iPtrStructPtr.val = 1
	// _ = toInt(&iPtrStructPtr)
	// _ = toIntPtrStructPtr(&iPtrStructPtr)
	// _ = toIntGeneric(&iPtrStructPtr)

	// calcs
	var c1i calcInt
	_ = sum(&c1i, 1, 2, 3)
	var c1g calcInt
	_ = sumGeneric(&c1g, 1, 2, 3)

	var c2i calcStruct
	_ = sum(&c2i, 1, 2, 3)
	var c2g calcStruct
	_ = sumGeneric(&c2g, 1, 2, 3)
}

func toInt(i Inter) int64 {
	return i.Int64()
}

func toInt32(i inter32) int64 {
	return i.Int64()
}

func toInt64(i inter64) int64 {
	return i.Int64()
}

func toIntStr(i interStr) int64 {
	return i.Int64()
}

func toIntPtr(i *interPtr) int64 {
	return i.Int64()
}

func toIntGeneric[T Inter](i T) int64 {
	return i.Int64()
}

func toIntStruct(i interStruct) int64 {
	return i.Int64()
}

func toIntPtrStruct(i *interPtrStruct) int64 {
	return i.Int64()
}

func toIntStructPtr(i interStructPtr) int64 {
	return i.Int64()
}

func toIntPtrStructPtr(i *interPtrStructPtr) int64 {
	return i.Int64()
}
