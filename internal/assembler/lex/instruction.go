// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/lex/instruction.go

package lex

import (
	"errors"
	"fmt"
	"github.com/codemicro/chip8/internal/assembler/token"
	"strconv"
	"strings"
)

func lexInstruction(peek func(offset int) rune, consume func() rune) (*token.Instruction, error) {

	var ins token.Instruction

	const (
		label = iota
		opcode
		operand
		comment
	)

	var buffer []rune
	var state int

lexLoop:
	for peek(0) != 0 {

		fmt.Println(string(buffer), state, string(peek(0)), peek(0), isWhitespace(peek(0)))

		switch state {
		case label:
			if isWhitespace(peek(0)) {
				consume()
				ins.Label = string(buffer)
				buffer = nil
				state = opcode
			} else if peek(0) == ':' && peek(1) == '\n' && isWhitespace(peek(2)) {
				consume()
				consume()
				consume()
				ins.Label = string(buffer)
				buffer = nil
				state = opcode
			} else {
				buffer = append(buffer, consume())
			}

		case opcode:
			opc, ok := lexOpcode(peek, consume)
			if !ok {
				return nil, errors.New("expecting opcode")
			}
			if x := peek(0); isWhitespace(x) {
				consume()
				ins.Opcode = strings.ToLower(opc)
				state = operand
			} else if x == '\n' || x == 0 {
				consume()
				ins.Opcode = strings.ToLower(opc)
				break lexLoop
			} else {
				return nil, fmt.Errorf("unexpected character %s, expecting newline or operand", string(peek(0)))
			}

		case operand:
			if x := peek(0); x == ';' || x == '\n' {
				consume()
				state = comment
			}

			if isWhitespace(peek(0)) {
				consume()
				continue lexLoop
			}

			if !(ins.Arg1 == nil || ins.Arg2 == nil) {
				state = comment
				break
			}

			opa, err := lexOperand(peek, consume)
			if err != nil {
				return nil, err
			}

			if ins.Arg1 == nil {
				ins.Arg1 = opa
			} else if ins.Arg2 == nil {
				ins.Arg2 = opa
			}

		case comment:
			break lexLoop
		}
	}

	return &ins, nil
}

func lexOpcode(peek func(offset int) rune, consume func() rune) (string, bool) {
	var o string
	for len(o) <= 4 {
		if x := peek(0); isWhitespace(x) || x == '\n' || x == 0 {
			if len(o) == 0 {
				return "", false
			}
			return o, true
		}
		o += string(consume())
	}
	return o, true
}

func lexOperand(peek func(offset int) rune, consume func() rune) (*token.Operand, error) {

	var buf []rune

	for {
		if x := peek(0); isWhitespace(x) || x == '\n' {
			break
		} else if x == 0 {
			return nil, errors.New("EOF when parsing operand")
		}
		buf = append(buf, consume())
	}

	instr := strings.ToLower(string(buf))

	base := 10
	opType := token.TypeValue
	if strings.HasPrefix(instr, "$") {
		// is a register
		instr = strings.TrimPrefix(instr, "$")
		base = 16
		opType = token.TypeRegister
	} else if strings.HasPrefix(instr, "0x") {
		instr = strings.TrimPrefix(instr, "0x")
		base = 16
	} else if strings.HasPrefix(instr, "0b") {
		instr = strings.TrimPrefix(instr, "0b")
		base = 2
	}

	n, err := strconv.ParseInt(instr, base, 32)
	if err != nil {
		return nil, err
	}

	return &token.Operand{
		OperandType: opType,
		Value:       int(n),
	}, nil
}
