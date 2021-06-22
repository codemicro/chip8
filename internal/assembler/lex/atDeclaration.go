// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/assembler/lex/atDeclaration.go

package lex

import (
	"errors"
	"fmt"
	"github.com/codemicro/chip8/internal/assembler/token"
	"strings"
)

const (
	keywordDefine = "define"
	keywordInclude = "include"
	keywordMacro = "macro"
	keywordSubroutine = "subroutine"
)

func lexAtDeclaration(peek func(offset int) rune, consume func() rune) (token.Token, error) {

	if peek(0) != '@' {
		return nil, errors.New("expecting @")
	}
	consume()

	if strings.EqualFold(peekMultiple(peek, 0, len(keywordDefine)), keywordDefine) {
		consumeMultiple(consume, len(keywordDefine))
		return lexDefine(peek, consume)
	} else if strings.EqualFold(peekMultiple(peek, 0, len(keywordInclude)), keywordInclude) {
		consumeMultiple(consume, len(keywordInclude))
		return lexInclude(peek, consume)
	} else if strings.EqualFold(peekMultiple(peek, 0, len(keywordMacro)), keywordMacro) {
		consumeMultiple(consume, len(keywordMacro))
		return lexMacro(peek, consume)
	} else if strings.EqualFold(peekMultiple(peek, 0, len(keywordSubroutine)), keywordSubroutine) {
		consumeMultiple(consume, len(keywordSubroutine))
		return lexSubroutine(peek, consume)
	}

	return nil, errors.New("unknown @ declaration")
}

func lexDefine(peek func(offset int) rune, consume func() rune) (*token.Define, error) {

	if peek(0) != ' ' {
		return nil, errors.New("expecting value in @define")
	}
	consume()

	// label
	label, err := lexLabel(peek, consume)
	if err != nil {
		return nil, err
	}

	if peek(0) != ' ' {
		return nil, errors.New("expecting value in @define")
	}
	consume()

	// value
	val, err := lexValue(peek, consume)
	if err != nil {
		return nil, err
	}
	if val.Type() != token.TypeValue {
		return nil, errors.New("define must define a constant value")
	}

	if x := peek(0); !(x == '\n' || x == 0) {
		return nil, errors.New("expected end of line")
	}

	return &token.Define{
		Label: label,
		Value: val,
	}, nil
}

func lexInclude(peek func(offset int) rune, consume func() rune) (*token.Include, error) {

	if peek(0) != ' ' {
		return nil, errors.New("expecting value in @include")
	}
	consume()

	var o []rune
	for {
		if peek(0) == '\n' || peek(0) == 0 {
			consume()
			break
		}
		o = append(o, consume())
	}

	return &token.Include{
		Filename: string(o),
	}, nil
}

func lexMacro(peek func(offset int) rune, consume func() rune) (*token.Macro, error) {

	const keywordEndMacro = "@endmacro"

	if peek(0) != ' ' {
		return nil, errors.New("expecting label for @macro")
	}
	consume()

	// lex label

	var label []rune

	for !(peek(0) == ':' || peek(0) == ' ') {
		if isValidIdentifier(peek(0)) {
			label = append(label, consume())
		} else if peek(0) == 0 {
			return nil, errors.New("unexpected EOF when parsing macro")
		} else {
			return nil, fmt.Errorf("disallowed character '%v' in label", string(peek(0)))
		}
	}
	consume()

	// lex arguments

	var args []*token.Argument

	if peek(-1) != ':' { // if we've not finished the `@macro label` line and have arguments to parse

		var buf []rune
		for {
			if isWhitespace(peek(0)) || peek(0) == ':' {

				isEnd := peek(0) == ':'

				consume()

				if len(buf) == 0 {
					return nil, errors.New("empty macro argument")
				}

				var argumentType token.Type
				isRegister := buf[0] == '&'
				if isRegister {
					argumentType = token.TypeRegister
					buf = buf[1:]
				} else {
					argumentType = token.TypeValue
				}

				args = append(args, &token.Argument{
					ArgumentType: argumentType,
					Label:        string(buf),
				})

				buf = nil

				if isEnd {
					break
				}

			} else if isValidIdentifier(peek(0)) || peek(0) == '&' {
				buf = append(buf, consume())
			} else {
				return nil, fmt.Errorf("disallowed character '%v' in label", string(peek(0)))
			}
		}

	}

	if peek(0) == '\n' {
		consume()
	}

	// lex instructions

	var instructions []*token.Instruction

	for {
		if peekMultiple(peek, 0, len(keywordEndMacro)) == keywordEndMacro {
			consumeMultiple(consume, len(keywordEndMacro))
			break
		}

		if peek(0) == 0 {
			return nil, errors.New("unexpected EOF while parsing macro")
		}

		ins, err := lexInstruction(peek, consume)
		if err != nil {
			return nil, err
		}
		fmt.Println(ins.Opcode)
		if ins.Label != "" {
			return nil, errors.New("macros cannot have labels in the macro body")
		}
		instructions = append(instructions, ins)

		if peek(0) == '\n' {
			consume()
		}
	}

	return &token.Macro{
		Instructions: instructions,
		Label:        string(label),
		Arguments:    args,
	}, nil
}

func lexSubroutine(peek func(offset int) rune, consume func() rune) (*token.Subroutine, error) {

	const keywordEndSubroutine = "@endsubroutine"

	// lex label

	if peek(0) != ' ' {
		return nil, errors.New("expecting label for @subroutine")
	}
	consume()

	// lex label

	var label []rune

	for peek(0) != ':' {
		fmt.Println(string(peek(0)))
		if isValidIdentifier(peek(0)) {
			label = append(label, consume())
		} else if peek(0) == 0 {
			return nil, errors.New("unexpected EOF when parsing subroutine")
		} else {
			return nil, fmt.Errorf("disallowed character '%v' in label", string(peek(0)))
		}
	}
	consume()

	if peek(0) == '\n' {
		consume()
	}

	// lex instructions

	var instructions []*token.Instruction

	for {
		fmt.Println(peekMultiple(peek, 0, len(keywordEndSubroutine)))

		if peekMultiple(peek, 0, len(keywordEndSubroutine)) == keywordEndSubroutine {
			consumeMultiple(consume, len(keywordEndSubroutine))
			break
		}

		if peek(0) == 0 {
			return nil, errors.New("unexpected EOF while parsing subroutine")
		}

		ins, err := lexInstruction(peek, consume)
		if err != nil {
			return nil, err
		}
		fmt.Println(ins.Opcode)
		if ins.Label != "" {
			return nil, errors.New("subroutines cannot have labels in the subroutine body")
		}
		instructions = append(instructions, ins)

		if peek(0) == '\n' {
			consume()
		}
	}

	return &token.Subroutine{
		Instructions: instructions,
		Label:        string(label),
	}, nil
}