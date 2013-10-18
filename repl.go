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

func main() {
	flag.Parse()
	if flag.NArg() >= 1 {
		data, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Println(err)
		}
		eval.EvalExpr(string(data))
	} else {
		fmt.Println("Welcome to Calc REPL", version)

		for {
			fmt.Print(">>>")
			in := bufio.NewReader(os.Stdin)
			str, _ := in.ReadString('\n')
			if str[:len(str)-2] == "q" {
				fmt.Println("QUIT!")
				break
			}
			res := eval.EvalExpr(str)
			if res != nil {
				fmt.Println(res)
			}
		}
	}
}
