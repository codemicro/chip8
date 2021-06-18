package lex

import (
	"fmt"
	"github.com/codemicro/chip8/internal/assembler/token"
	"testing"
)


func Test_lexAtDeclaration(t *testing.T) {
	input := []byte(`@macro banana a bc &dfskjgh:
	hello $f
	hi 647
	words
@endmacro`)

	inputLength := len(input)
	var index int

	peek := func(offset int) rune {
		if index+offset >= inputLength {
			return 0
		}
		// fmt.Printf("%#v\n", string(input[index+offset:]))
		return rune(input[index+offset])
	}

	consume := func() rune {
		if index >= inputLength {
			return 0
		}
		index += 1
		return rune(input[index-1])
	}

	ins, err := lexAtDeclaration(peek, consume)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", ins)
		fmt.Println("  ARGS")
		for _, arg := range ins.(*token.Macro).Arguments {
			fmt.Printf("  %#v\n", arg)
		}

		fmt.Println("  INS'")
		for _, arg := range ins.(*token.Macro).Statements {
			fmt.Printf("  %#v - %s\n", arg, arg)
		}
	}
}