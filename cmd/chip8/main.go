package main

import (
	"github.com/codemicro/chip8/internal/display"
	vm2 "github.com/codemicro/chip8/internal/vm"
)

func main() {
	disp := display.NewDisplay(64, 32, 5)

	vm := vm2.NewChip8(disp)
	go vm.Run()

	if err := disp.Start(); err != nil {
		panic(err)
	}
}
