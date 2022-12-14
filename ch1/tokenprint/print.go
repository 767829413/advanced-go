package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
)

func main() {
	// src is the input that we want to tokenize.
	src, _ := os.ReadFile(`./tmp/main.go`)

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	fmt.Printf("%s\t%s\t%s\n", "pos", "token", "literal")
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}
}
