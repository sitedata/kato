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

package event

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEvent(t *testing.T) {
	i := []int{}
	f := func(s interface{}) {
		i = append(i, 1)
	}
	f2 := func(s interface{}) {
		i = append(i, 2)
		i = append(i, 3)
	}

	Convey("events package test", t, func() {
		Convey("init events package should be a success", func() {
			So(len(i), ShouldEqual, 0)
			So(len(Events[EXIT]), ShouldEqual, 0)
		})

		Convey("empty events execute Off function should not be a success", func() {
			So(Off(EXIT, f), ShouldNotBeNil)
		})

		Convey("multi execute On function for a function should not be a success", func() {
			So(On(EXIT, f), ShouldBeNil)
			So(On(EXIT, f), ShouldNotBeNil)
		})

		Convey("execute Emit function should be a success", func() {
			Emit(EXIT, nil)
			So(len(i), ShouldEqual, 1)
		})

		Convey("events package should be working", func() {
			So(On(EXIT, f2), ShouldBeNil)
			So(len(Events[EXIT]), ShouldEqual, 2)
			So(len(i), ShouldEqual, 1)

			So(Off(EXIT, f), ShouldBeNil)
			So(len(Events[EXIT]), ShouldEqual, 1)
		})
	})
}
