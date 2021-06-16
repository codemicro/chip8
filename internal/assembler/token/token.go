// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/token/token.go

package token

import "fmt"

type Type uint8

const (
	TypeDefine Type = iota
	TypeInclude
	TypeMacro
	TypeSubroutine
	TypeInstruction

	TypeRegister
	TypeValue
)

type Token interface {
	fmt.Stringer
	Type() Type
}

type Operand struct {
	OperandType Type
	Value       int
}

func (o *Operand) Type() Type { return o.OperandType }
func (o *Operand) String() string {

	if o == nil {
		return ""
	}

	var chr string
	if o.Type() == TypeRegister {
		chr = "$"
	}
	return fmt.Sprintf("%s%d", chr, o.Value)
}

type Instruction struct {
	Label  string
	Opcode string
	Arg1   *Operand
	Arg2   *Operand
}

func (i *Instruction) Type() Type { return TypeInstruction }
func (i *Instruction) String() string {
	return fmt.Sprintf(
		"%s %s %s %s",
		i.Label,
		i.Opcode,
		i.Arg1.String(),
		i.Arg2.String(),
	)
}
