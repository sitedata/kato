// Copyright (C) 2021 Gridworkz Co., Ltd.
// KATO, Application Management Platform

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

package utils

import (
	"math/rand"
	"time"
)

// ASCII values 33 ~ 126
const _dcl = 126 - 33 + 1

var defaultCharacters [_dcl]byte

func init() {
	for i := 0; i < _dcl; i++ {
		defaultCharacters[i] = byte(i + 33)
	}

	rand.Seed(time.Now().UnixNano())
}

func RandString(length int, characters ...byte) string {
	if len(characters) == 0 {
		characters = defaultCharacters[:]
	}

	n := len(characters)
	var rs = make([]byte, length)

	for i := 0; i < length; i++ {
		rs[i] = characters[rand.Intn(n-1)]
	}

	return string(rs)
}
