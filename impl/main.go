package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	replFlag = flag.Bool("repl", false, "Runs the funlang REPL")
	astFlag  = flag.Bool("ast", false, "Prints the AST to stderr. (debug)")
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		*replFlag = true
	}

	ev := Eval{defs: map[string]Node{}}
	code := []Node{}

	for i := 0; i < len(args); i++ {
		file, err := os.Open(args[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "File %s doesn't exist.\n", args[i])
			return
		}

		b, err := io.ReadAll(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read %s's content.\n", args[i])
			return
		}

		par := Parser{
			lexer: Lexer{buf: []rune(string(b))},
			defs:  ev.defs,
		}

		ast, err := par.topLevel()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: error: %v\n", err)
		}

		if *astFlag {
			fmt.Fprintf(os.Stderr, "%s\n", ast.toString(0))
		}

		ev.defs = par.defs

		code = append(code, ast)
	}

	for i := 0; i < len(code); i++ {
		_, err := ev.expression(code[i], 0, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s runtime error: %v", args[i], err)
		}
	}

	if *replFlag {
		runRepl(ev)
	}
}
