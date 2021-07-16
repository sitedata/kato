// KATO, Application Management Platform
// Copyright (C) 2021 Gridworkz Co., Ltd.

// Permission is hereby granted, free of charge, to any person obtaining a copy of this 
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons 
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies or 
// substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, 
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
// PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
// FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
// ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package util

import (
	"reflect"
	"unsafe"
)

//ToString []byte to string
//BenchmarkToString-4     30000000                42.0 ns/op
func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

//ToByte string to []byte
//BenchmarkToByte-4       2000000000               0.36 ns/op
func ToByte(v string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return *(*[]byte)(unsafe.Pointer(sh))
}

//StringArrayContains string array contains
func StringArrayContains(list []string, source string) bool {
	for _, l := range list {
		if l == source {
			return true
		}
	}
	return false
}

//Reverse reverse sort string array
func Reverse(source []string) []string {
	for i := 0; i < len(source)/2; i++ {
		source[i], source[len(source)-i-1] = source[len(source)-i-1], source[i]
	}
	return source
}
