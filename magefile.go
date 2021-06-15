// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: magefile.go

//+build mage

package main

import (
	"fmt"
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