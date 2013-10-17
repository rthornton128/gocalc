package main

import (
	"bufio"
	"fmt"
	"misc/calc/eval"
	"os"
)

var version = "0.1"

func main() {
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
