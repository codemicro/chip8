// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/lex/instruction_test.go

package lex

import (
	"fmt"
	"testing"
)

func Test_lexInstruction(t *testing.T) {
	input := []byte("main:\n blah $2 1 ; hi there!")

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

	ins, err := lexInstruction(peek, consume)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ins.String())
	}

}
