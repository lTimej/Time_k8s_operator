package utils

import "unsafe"

/*
pointer是任意类型的指针，可以指向任意类型数据s
*/

func String2Bytes(str string) []byte {
	return *(*[]byte)(unsafe.Pointer(&struct {
		string
		int
	}{str, len(str)}))
}

func Bytes2String(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}
