package main

import (
	"bufio"
	"fmt"
	"monkey/lexing"
	"os"
)

const PROMPT = ">> "

func main() {
	fmt.Println("Monkey language RLPL (Read Lex Print Loop).")

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

		token := &lexing.Token{Type: lexing.TOKEN_ILLEGAL}
		for token.Type != lexing.TOKEN_EOF {
			token = lexer.Next()
			fmt.Printf(
				"%s: %q\n",
				lexing.TokenTypeToString(token.Type),
				token.Literal,
			)
		}
	}

}
