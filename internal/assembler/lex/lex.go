// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/lex/lex.go

package lex

import (
	"fmt"
	"github.com/codemicro/chip8/internal/assembler/token"
)

func Lex(input []byte) ([]token.Token, error) {

	inputLength := len(input)
	var index int

	peek := func(offset int) rune {
		if index+offset >= inputLength {
			return 0
		}
		return rune(input[index+offset])
	}

	consume := func() rune {
		if index >= inputLength {
			return 0
		}
		index += 1
		return rune(input[index-1])
	}

	var tokens []token.Token

	for index < inputLength {
		if peekMultiple(peek, 0, 2) == "\n\n" {
			fmt.Println("Hi!")
			consumeMultiple(consume, 2)
		}
		if peek(0) == '@' {
			tk, err := lexAtDeclaration(peek, consume)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tk)
		} else if peek(0) == ';' {
			for peek(0) != '\n' {
				consume()
			}
			consume()
		} else {
			tk, err := lexInstruction(peek, consume)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tk)
		}
	}

	return tokens, nil
}

func lexLabel(peek func(offset int) rune, consume func() rune) (string, error) {
	var b []rune
	for {
		if isWhitespace(peek(0)) {
			break
		} else if isValidIdentifier(peek(0)) {
			b = append(b, consume())
		} else {
			return "", fmt.Errorf("disallowed character %#v in label", string(peek(0)))
		}
	}
	return string(b), nil
}

func isValidIdentifier(r rune) bool {
	return isDigit(r) || isCharacter(r)
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isCharacter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

func peekMultiple(peek func(offset int) rune, offset, runLength int) string {
	var o []rune
	for i := 0; i < runLength; i += 1 {
		o = append(o, peek(offset + i))
	}
	return string(o)
}

func consumeMultiple(consume func() rune, runLength int) string {
	var o []rune
	for i := 0; i < runLength; i += 1 {
		o = append(o, consume())
	}
	return string(o)
}