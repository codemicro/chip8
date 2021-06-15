// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: magefile.go

//+build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"os"
	"path/filepath"

	"github.com/codemicro/alib-go/mage/exmg"
	"github.com/magefile/mage/sh"
)

func Build() error {
	const buildPackage = "github.com/codemicro/chip8/cmd/chip8"

	outputDir := filepath.Join("bin", fmt.Sprintf("%s-%s", exmg.GetTargetOS(), exmg.GetTargetArch()))
	basePackageName := filepath.Base(buildPackage)

	_ = os.MkdirAll(outputDir, os.ModeDir)

	return sh.Run("go", "build", "-o", filepath.Join(outputDir, basePackageName), buildPackage)
}

func Test() error {
	const testPackage = "github.com/codemicro/chip8/..."

	var command = []string{"test", "-v"}
	var runFunc = sh.Run
	if mg.Verbose() {
		command = append(command, "-v")
		runFunc = sh.RunV
	}
	command = append(command, testPackage)

	return runFunc("go", command...)
}