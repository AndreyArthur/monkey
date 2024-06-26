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

func main() {
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
		object := evaluating.Eval(env, ast)

		fmt.Println(object.Inspect())
	}

}
