package main

import (
	"bufio"
	"fmt"
	"monkey/evaluating"
	"monkey/lexing"
	"monkey/parsing"
	"os"
	"strings"
)

const PROMPT = ">> "

func repl() {
	fmt.Println("Monkey language REPL (Read Eval Print Loop).")

	in := os.Stdin
	out := os.Stdout

	scanner := bufio.NewScanner(in)
	content := ""

	for {
		fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		current := scanner.Text()
		current = strings.TrimSpace(current)

		if current == "exit" {
			return
		}

		if current[len(current)-1] != ';' {
			current += ";"
		}

		lexer := lexing.NewLexer(content + current)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()

		if parser.HasErrors() {
			for _, error := range parser.GetErrors() {
				fmt.Println(error)
			}
			continue
		}

		content += current

		env := evaluating.NewEnvironment(nil)
		evaluating.InjectBuiltinFunctions(env)
		object := evaluating.Eval(env, ast)

		fmt.Println(object.Inspect())
	}

}

func file() {
	filepath := os.Args[1]

	content, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	lexer := lexing.NewLexer(string(content))
	parser := parsing.NewParser(lexer)
	ast := parser.Parse()

	if parser.HasErrors() {
		for _, error := range parser.GetErrors() {
			fmt.Println(error)
		}
		return
	}

	env := evaluating.NewEnvironment(nil)
	evaluating.InjectBuiltinFunctions(env)
	_ = evaluating.Eval(env, ast)
}

func main() {
	if len(os.Args) == 1 {
		repl()
	} else {
		file()
	}
}
