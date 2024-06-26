package main

import (
	"bufio"
	"fmt"
	"monkey/lexing"
	"monkey/parsing"
	"os"
)

const PROMPT = ">> "

func main() {
	fmt.Println("Monkey language RPPL (Read Parse Print Loop).")

	in := os.Stdin
	out := os.Stdout

	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		content := scanner.Text()

		if content == "exit" {
			return
		}

		lexer := lexing.NewLexer(content)
		parser := parsing.NewParser(lexer)
		ast := parser.Parse()

		if parser.HasErrors() {
			for _, error := range parser.GetErrors() {
				fmt.Println(error)
			}
		} else {
			fmt.Println(ast.String())
		}
	}

}
