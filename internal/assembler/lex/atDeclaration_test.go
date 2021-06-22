// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/lex/atDeclaration_test.go

package lex

import (
	"fmt"
	"testing"
)


func Test_lexAtDeclaration(t *testing.T) {
	input := []byte(`main:
    set $0 3
    set $1 3

    set $2 0xA
    char $2

    disp $0 $1 5`)

	ins, err := Lex(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", ins)
	}
}