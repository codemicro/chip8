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

