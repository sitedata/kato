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

import "testing"

func TestCheckIP(t *testing.T) {
	t.Logf("ip %s %t", "1829.123", CheckIP("1829.123"))
	t.Logf("ip %s %t", "1829.123.1", CheckIP("1829.123.1"))
	t.Logf("ip %s %t", "1829.123.2.1", CheckIP("1829.123.2.1"))
	t.Logf("ip %s %t", "y.123", CheckIP("y.123"))
	t.Logf("ip %s %t", "0.0.0.0", CheckIP("0.0.0.0"))
	t.Logf("ip %s %t", "127.0.0.1", CheckIP("127.0.0.1"))
	t.Logf("ip %s %t", "localhost", CheckIP("localhost"))
	t.Logf("ip %s %t", "192.168.0.1", CheckIP("192.168.0.1"))
}
