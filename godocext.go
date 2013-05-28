// Copyright 2013 Fredy Wijaya
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
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	helpFlag         *bool
	methodOnlyFlag   *bool
	functionOnlyFlag *bool
	typeOnlyFlag     *bool
)

func printUsage() {
	fmt.Println("Usage:", os.Args[0])
	flag.PrintDefaults()
}

func isMethod(bytes []byte) bool {
	r := regexp.MustCompile("^func \\(.*\\) .*\\(.*\\)")
	if r.Find(bytes) == nil {
		return false
	}
	return true
}

func isType(bytes []byte) bool {
	r := regexp.MustCompile("^type")
	if r.Find(bytes) == nil {
		return false
	}
	return true
}

func isFunction(bytes []byte) bool {
	r := regexp.MustCompile("^func .*\\(.*\\)")
	if r.Find(bytes) == nil {
		return false
	}
	return true
}

func formatGodocExecutable() string {
	godoc := "godoc"
	if runtime.GOOS == "windows" {
		return godoc + ".exe"
	}
	return godoc
}

func runGoDoc(packageName string) error {
	var godocPath string
	if path, e := exec.LookPath(formatGodocExecutable()); e != nil {
		goroot := os.Getenv("GOROOT")
		godocPath = filepath.Join(goroot, "bin", formatGodocExecutable())
	} else {
		godocPath = path
	}
	cmd := exec.Command(godocPath, packageName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if e := cmd.Start(); e != nil {
		return e
	}
	go func() {
		stdoutReader := bufio.NewReader(stdout)
		for true {
			bytes, _, e := stdoutReader.ReadLine()
			if e != nil {
				break
			}
			line := string(bytes)
			if *methodOnlyFlag {
				if !isMethod(bytes) {
					continue
				}
			}
			if *functionOnlyFlag {
				if isMethod(bytes) || !isFunction(bytes) {
					continue
				}
			}
			if *typeOnlyFlag {
				if !isType(bytes) {
					continue
				}
				if strings.HasSuffix(line, "{") {
					// to beautify the output
					line += " ... }"
				}
			}
			fmt.Println(packageName+":", line)
		}
	}()
	cmd.Wait()
	return nil
}

func init() {
	helpFlag = flag.Bool("h", false, "show help")
	functionOnlyFlag = flag.Bool("f", false, "show functions only")
	methodOnlyFlag = flag.Bool("m", false, "show methods only")
	typeOnlyFlag = flag.Bool("t", false, "show types only")

	flag.Parse()

	if *helpFlag {
		printUsage()
		os.Exit(0)
	}
}

func main() {
	goroot := os.Getenv("GOROOT")
	gosrc := filepath.Join(goroot, "src", "pkg")
	filepath.Walk(gosrc, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || gosrc == path {
			return nil
		}
		packageName := path[len(gosrc)+1:]
		runGoDoc(packageName)
		return nil
	})
}
