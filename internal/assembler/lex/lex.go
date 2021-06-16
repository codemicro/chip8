// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/lex/lex.go

package lex

import "go/token"

func Lex(input []byte) ([]token.Token, error) {

	//inputLength := len(input)
	//var index, status int
	//
	//peek := func(offset int) rune {
	//	if index+offset >= inputLength {
	//		return 0
	//	}
	//	return rune(input[index+offset])
	//}
	//
	//consume := func() rune {
	//	if index >= inputLength {
	//		return 0
	//	}
	//	index += 1
	//	return rune(input[index-1])
	//}
	//
	//const (
	//	Instruction = iota
	//)
	//
	//for index < inputLength {
	//
	//	switch status {
	//	case Instruction:
	//
	//	}
	//
	//}

	return nil, nil
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}