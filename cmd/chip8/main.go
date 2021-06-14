package main

import (
	"github.com/alexflint/go-arg"
	"github.com/codemicro/chip8/internal/ui"
	vm2 "github.com/codemicro/chip8/internal/vm"
	"io/ioutil"
	"path/filepath"
)

var args struct {
	InputFile string `arg:"positional"`
	DebugMode bool   `arg:"-d,-v,--verbose" help:"enable verbose/debug mode"`
	UIScale   int    `arg:"-s,--scale" help:"UI scale factor" default:"5"`
	ToneFrequency int `arg:"--frequency" help:"sound timer tone frequency" default:"350"`
	ClockSpeed int `arg:"-c,--clock" help:"approximate clock speed in hertz" default:"500"`
}

func main() {

	arg.MustParse(&args)

	fcont, err := ioutil.ReadFile(args.InputFile)
	if err != nil {
		panic(err)
	}

	disp, err := ui.NewUI(64, 32, args.UIScale, filepath.Base(args.InputFile), args.ToneFrequency)
	if err != nil {
		panic(err)
	}
	// disp.Debug = true

	vm := vm2.NewChip8(fcont, disp, args.ClockSpeed)
	vm.Debug = args.DebugMode
	go vm.Run()

	if err = disp.Start(); err != nil {
		panic(err)
	}
}
