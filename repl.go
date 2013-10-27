package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"misc/calc/eval"
	"os"
)

var version = "0.1"

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

func main() {
	flag.Parse()
	if flag.NArg() >= 1 {
		data, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Println(err)
		}
		eval.EvalExpr(string(stripCR(data)))
	} else {
		fmt.Println("Welcome to Calc REPL", version)

		for {
			fmt.Print(">>>")
			in := bufio.NewReader(os.Stdin)
			b, _ := in.ReadBytes('\n')
			b = stripCR(b)
			if len(b) <= 2 && b[0] == 'q' {
				fmt.Println("QUIT!")
				break
			}
			res := eval.EvalExpr(string(b))
			if res != nil {
				fmt.Println(res)
			}
		}
	}
}
