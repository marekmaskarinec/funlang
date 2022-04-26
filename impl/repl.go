package main

import (
	"bufio"
	"fmt"
	"os"
)

func runRepl(ev Eval) {
	var result Value

	fmt.Println("funlang repl. Write \"bye\" or press Ctrl+D to exit.")

	for {
		fmt.Print(">>> ")
		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}
		cmd = cmd[:len(cmd)-1]

		if cmd == "bye" {
			break
		}

		par := Parser{
			lexer: Lexer{buf: []rune(cmd)},
			defs:  ev.defs,
		}
		ast, err := par.topLevel()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			continue
		}

		ev.defs = par.defs
		tmpResult, err := ev.expression(ast, result, result)
		if err != nil {
			fmt.Printf("runtime error: %v\n", err)
			continue
		} else {
			result = tmpResult
			fmt.Println("result:", result)
		}
	}
}
