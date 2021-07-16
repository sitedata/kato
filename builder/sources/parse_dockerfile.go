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

package sources

import (
	"io"
	"os"
	"sort"

	"github.com/gridworkz/kato/util/dockerfile/command"
	"github.com/gridworkz/kato/util/dockerfile/parser"
)

// Command Represents a single line (layer) in a Dockerfile.
// For example `FROM ubuntu:xenial`
type Command struct {
	Cmd       string   // lowercased command name (ex: `from`)
	SubCmd    string   // for ONBUILD only this holds the sub-command
	Json      bool     // whether the value is written in json form
	Original  string   // The original source line
	StartLine int      // The original source line number
	Flags     []string // Any flags such as `--from=...` for `COPY`.
	Value     []string // The contents of the command (ex: `ubuntu:xenial`)
}

// IOError A failure in opening a file for reading.
type IOError struct {
	Msg string
}

func (e IOError) Error() string {
	return e.Msg
}

// ParseError A failure in parsing the file as a dockerfile.
type ParseError struct {
	Msg string
}

func (e ParseError) Error() string {
	return e.Msg
}

// AllCmds List all legal cmds in a dockerfile
func AllCmds() []string {
	var ret []string
	for k := range command.Commands {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

// ParseReader Parse a Dockerfile from a reader.  A ParseError may occur.
func ParseReader(file io.Reader) ([]Command, error) {
	directive := parser.Directive{LookingForDirectives: true}
	parser.SetEscapeToken(parser.DefaultEscapeToken, &directive)
	ast, err := parser.Parse(file, &directive)
	if err != nil {
		return nil, ParseError{err.Error()}
	}

	var ret []Command
	for _, child := range ast.Children {
		cmd := Command{
			Cmd:       child.Value,
			Original:  child.Original,
			StartLine: child.StartLine,
			Flags:     child.Flags,
		}

		// Only happens for ONBUILD
		if child.Next != nil && len(child.Next.Children) > 0 {
			cmd.SubCmd = child.Next.Children[0].Value
			child = child.Next.Children[0]
		}

		cmd.Json = child.Attributes["json"]
		for n := child.Next; n != nil; n = n.Next {
			cmd.Value = append(cmd.Value, n.Value)
		}

		ret = append(ret, cmd)
	}
	return ret, nil
}

// ParseFile Parse a Dockerfile from a filename.  An IOError or ParseError may occur.
func ParseFile(filename string) ([]Command, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, IOError{err.Error()}
	}
	defer file.Close()

	return ParseReader(file)
}
