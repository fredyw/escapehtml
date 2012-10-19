// Copyright 2012 Fredy Wijaya
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
    "errors"
    "fmt"
    "html"
    "io/ioutil"
    "os"
    "path/filepath"
    "strings"
)

func printUsage() {
    fmt.Println("Usage:", os.Args[0],
        "<source_file/source_dir>", "[dest_dir]")
}

func errorMessage(s string) string {
    return "Error: " + s
}

type fileType struct {
    directory   bool
    regularFile bool
}

func fileExists(path string) (*fileType, error) {
    file, err := os.Open(path)
    if err != nil {
        if os.IsNotExist(err) {
            return nil, errors.New(
                errorMessage(path + " does not exist"))
        }
    }
    defer file.Close()
    fileInfo, err := file.Stat()
    if err != nil {
        return nil, err
    }
    if fileInfo.IsDir() {
        return &fileType{directory: true}, nil
    }
    return &fileType{regularFile: true}, nil
}

func validateArgs() (bool, error) {
    if len(os.Args) != 2 && len(os.Args) != 3 {
        return false, nil
    }

    if _, err := fileExists(os.Args[1]); err != nil {
        return false, err
    }
    if len(os.Args) == 3 {
        ft, _ := fileExists(os.Args[2])
        if ft != nil && !ft.directory {
            return false, errors.New(
                errorMessage(os.Args[2] + " must be a directory"))
        }
    }
    return true, nil
}

func printHeader(header string) {
    fmt.Println(strings.Repeat("=", 72))
    fmt.Println(header)
    fmt.Println(strings.Repeat("=", 72))
}

func escapeHTML(srcPath string, destPath string) {
    filepath.Walk(srcPath,
        func(path string, info os.FileInfo, err error) error {
            if info.IsDir() {
                return nil
            }

            b, err := ioutil.ReadFile(path)
            if err != nil {
                return nil
            }

            escapedString := html.EscapeString(string(b))
            if destPath == "" {
                printHeader(path)
                fmt.Println(escapedString)
            } else {
                if ft, _ := fileExists(destPath); ft == nil {
                    if e := os.MkdirAll(destPath, 0775); e != nil {
                        errors.New("Unable to create directory: " +
                            destPath)
                    } else {
                        fmt.Println("Creating directory: " + destPath)
                    }
                }
                newPath := filepath.Join(destPath,
                    filepath.Base(path)+".txt")
                fmt.Println("Creating", newPath)
                e := ioutil.WriteFile(newPath, []byte(escapedString), 0644)
                if e != nil {
                    return e
                }
            }
            return nil
        })
}

func main() {
    if valid, err := validateArgs(); !valid {
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        } else {
            printUsage()
            os.Exit(1)
        }
    }

    srcPath := os.Args[1]
    destPath := ""
    if len(os.Args) == 3 {
        destPath = os.Args[2]
    }
    escapeHTML(srcPath, destPath)
}
