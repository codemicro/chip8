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

		switch state {
		case label:
			if isWhitespace(peek(0)) {
				consume()
				ins.Label = string(buffer)
				buffer = nil
				state = opcode
			} else if peek(0) == ':' && peek(1) == '\n' && isWhitespace(peek(2)) { // label, colon, newline format
				consume()
				consume()
				consume()
				ins.Label = string(buffer)
				buffer = nil
				state = opcode
			} else if !isValidIdentifier(peek(0)) {
				return nil, fmt.Errorf("disallowed character %#v in label", string(peek(0)))
			} else {
				buffer = append(buffer, consume())
			}
		case opcode:
			opc, err := lexOpcode(peek, consume)
			if err != nil {
				return nil, err
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
				return nil, fmt.Errorf("unexpected character %#v, expecting newline or operand", string(peek(0)))
			}

		case operand:
			if x := peek(0); x == ';' {
				consume()
				state = comment
			} else if  x == '\n' {
				consume()
				break lexLoop
			}

			if isWhitespace(peek(0)) {
				consume()
				continue lexLoop
			}

			if !(ins.Arg1 == nil || ins.Arg2 == nil) {
				state = comment
				break
			}

			opa, err := lexValue(peek, consume)
			if err != nil {
				return nil, err
			}

			if ins.Arg1 == nil {
				ins.Arg1 = opa
			} else if ins.Arg2 == nil {
				ins.Arg2 = opa
			}

		case comment:
			if peek(0) == '\n' || peek(0) == 0 {
				consume()
				break lexLoop
			}
			consume()
		}
	}

	return &ins, nil
}

func lexOpcode(peek func(offset int) rune, consume func() rune) (string, error) {
	var o string
	for len(o) <= 4 {
		if x := peek(0); isWhitespace(x) || x == '\n' || x == 0 {
			if len(o) == 0 {
				return "", errors.New("expecting opcode")
			}
			return o, nil
		}
		o += string(consume())
	}
	return o, nil
}

func lexValue(peek func(offset int) rune, consume func() rune) (*token.Operand, error) {

	var buf []rune

	for {
		if x := peek(0); isWhitespace(x) || x == '\n' || (x == 0 && len(buf) != 0) {
			break
		} else if x == 0 {
			return nil, errors.New("EOF when parsing value")
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
