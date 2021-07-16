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
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"os"
)

const filechunk = 8192 // we settle for 8KB
//CreateFileHash compute sourcefile hash and write hashfile
func CreateFileHash(sourceFile, hashfile string) error {
	file, err := os.Open(sourceFile)
	if err != nil {
		return err
	}
	defer file.Close()
	fileinfo, _ := file.Stat()
	if fileinfo.IsDir() {
		return fmt.Errorf("do not support compute folder hash")
	}
	writehashfile, err := os.OpenFile(hashfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0655)
	if err != nil {
		return fmt.Errorf("create hash file error %s", err.Error())
	}
	defer writehashfile.Close()
	if fileinfo.Size() < filechunk {
		return createSmallFileHash(file, writehashfile)
	}
	return createBigFileHash(file, writehashfile)
}

func createBigFileHash(sourceFile, hashfile *os.File) error {
	// calculate the file size
	info, _ := sourceFile.Stat()
	filesize := info.Size()
	blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))
	hash := md5.New()

	for i := uint64(0); i < blocks; i++ {
		blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
		buf := make([]byte, blocksize)
		index, err := sourceFile.Read(buf)
		if err != nil {
			return err
		}
		// append into the hash
		_, err = hash.Write(buf[:index])
		if err != nil {
			return err
		}
	}
	_, err := hashfile.Write([]byte(fmt.Sprintf("%x", hash.Sum(nil))))
	if err != nil {
		return err
	}
	return nil
}

func createSmallFileHash(sourceFile, hashfile *os.File) error {
	md5h := md5.New()
	_, err := io.Copy(md5h, sourceFile)
	if err != nil {
		return err
	}
	_, err = hashfile.Write([]byte(fmt.Sprintf("%x", md5h.Sum(nil))))
	if err != nil {
		return err
	}
	return nil
}

//CreateHashString
func CreateHashString(source string) (hashstr string, err error) {
	md5h := md5.New()
	_, err = md5h.Write([]byte(source))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5h.Sum(nil)), nil
}
