// Copyright (c) 2013, Rob Thornton
// All rights reserved.
// This software is governed by a Simplied BSD-License. Please see the
// LICENSE included in this distribution for a copy of the full license
// or, if one is not included, you may also find a copy at
// http://opensource.org/licenses/BSD-2-Clause

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/rthornton128/gocalc/eval"
	"github.com/rthornton128/gocalc/token"
	"github.com/rthornton128/gocalc/trans"
	"io/ioutil"
	"os"
	"path/filepath"
)

var version = "0.2"

func stripCR(in []byte) []byte {
	out := make([]byte, len(in))
	i := 0
	for _, ch := range in {
		if ch != '\r' {
			out[i] = ch
			i++
		}
	}
	return out[:i]
}

func transpileCode(out *os.File, expr string) {
  /*data, err := ioutil.ReadFile(flag.Arg(0))
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }*/
  trans.TransExpr(out, string(stripCR([]byte(expr))))
}


func main() {
  t := flag.Bool("t", false, "Transpile")
	flag.Parse()
	if flag.NArg() >= 1 {

    if (*t == true) {
      fmt.Println("transpiling")
      out, err := os.Create("out.c")
      if err != nil {
        fmt.Println("create:", err)
        os.Exit(1)
      }
			data, err := ioutil.ReadFile(flag.Arg(0))
			if err != nil {
				fmt.Println(err)
			}
      transpileCode(out, string(data))//flag.Arg(0))
      return
    }


		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer f.Close()

		var fi os.FileInfo
		fi, err = f.Stat()
		if !fi.IsDir() {
			data, err := ioutil.ReadFile(flag.Arg(0))
			if err != nil {
				fmt.Println(err)
			} else {
				eval.EvalFile(flag.Arg(0), string(stripCR(data)))
			}
		} else {
			fset := token.NewFileSet()
			names, err := f.Readdirnames(0)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			for _, name := range names {
				data, err := ioutil.ReadFile(filepath.Join(flag.Arg(0), name))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				fset.AddFile(name, string(stripCR(data)))
			}
			eval.EvalPackage(flag.Arg(0), fset)
		}
	} else {
		fmt.Println("Welcome to Calc REPL", version)
		fmt.Println()
		fmt.Println("Type in expression(s) to evaluate on one or more lines.")
		fmt.Println("Press enter on an empty line to execute the expression(s).")
		fmt.Println("Type 'q' (without quotes) on an empty line to exit.")

		in := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(">>>")
			var expr string
			stop := false
			for !stop {
				b, _ := in.ReadBytes('\n')
				b = stripCR(b)
				if len(b) <= 2 {
					switch b[0] {
					case 'q':
						fmt.Println("QUIT!")
						os.Exit(0)
					case '\n':
						stop = true
					}
				}
				expr += string(b)
			}
			res := eval.EvalExpr(expr)
			if res != nil {
				fmt.Println(res)
			}
		}
	}
}
