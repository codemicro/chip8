package main

import (
	"github.com/codemicro/chip8/internal/display"
	vm2 "github.com/codemicro/chip8/internal/vm"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	filename := os.Args[1]
	fcont, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	disp := display.NewDisplay(64, 32, 5, filepath.Base(filename))
	//disp.Debug = true

	vm := vm2.NewChip8(fcont, disp, 700)
	// vm.Debug = true
	go vm.Run()

	if err = disp.Start(); err != nil {
		panic(err)
	}
}
