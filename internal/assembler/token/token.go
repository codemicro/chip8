// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/token/token.go

package token

import (
	"fmt"
	"strings"
)

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
	Arg1   *Operand // Arg1 may be nil
	Arg2   *Operand // Arg2 may be nil
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

type Define struct {
	Label string
	Value *Operand // Value may not be nil
}

func (d *Define) Type() Type { return TypeDefine }
func (d *Define) String() string { return fmt.Sprintf("define %s as %s", d.Label, d.Value.String()) }

type Include struct {
	Filename string
}

func (i *Include) Type() Type { return TypeInclude }
func (i *Include) String() string { return fmt.Sprintf("include %s", i.Filename) }

type Argument struct {
	ArgumentType Type
	Label        string
}

func (a *Argument) String() string {

	if a == nil {
		return ""
	}

	var chr string
	if a.ArgumentType == TypeRegister {
		chr = "$"
	}
	return fmt.Sprintf("%s%s", chr, a.Label)
}

type Macro struct {
	Instructions []*Instruction
	Label        string
	Arguments    []*Argument // Arguments may not have nil values
}

func (m *Macro) Type() Type { return TypeMacro }
func (m *Macro) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("macro %s %v\n", m.Label, m.Arguments))

	for _, ins := range m.Instructions {
		sb.WriteString("  ")
		sb.WriteString(ins.String())
	}

	return sb.String()
}

type Subroutine struct {
	Label string
	Instructions []*Instruction
}

func (s *Subroutine) Type() Type { return TypeSubroutine }
func (s *Subroutine) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("macro %s\n", s.Label))

	for _, ins := range s.Instructions {
		sb.WriteString("  ")
		sb.WriteString(ins.String())
	}

	return sb.String()
}