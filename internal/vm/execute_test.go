// https://github.com/codemicro/chip8
// Copyright (c) 2021, codemicro and contributors
// SPDX-License-Identifier: MIT
// Filename: internal/vm/execute_test.go

package vm

import "testing"

func vmFixture(instruction []byte) (*Chip8, *uid) {
	c, u := vmFixtureWithoutTick(instruction)

	if instruction != nil {
		c.tick()
	}

	return c, u
}

func vmFixtureWithoutTick(instruction []byte) (*Chip8, *uid) {
	u := &uid{}
	c := NewChip8(instruction, u, 500)

	return c, u
}

type uid struct {
	output *[32][64]bool
	keys [][]uint8
}

func (u *uid) PublishNewDisplay(in [32][64]bool) {
	u.output = &in
}
func (u *uid) GetPressedKeys() []uint8 {
	if len(u.keys) > 0 {
		x := u.keys[0]
		u.keys = u.keys[1:]
		return x
	}
	return nil
}
func (u *uid) StartTone() {}
func (u *uid) StopTone() {}

func Test_ClearScreen(t *testing.T) {
	c, u := vmFixtureWithoutTick(nil)

	c.clearScreen()

	if u.output == nil {
		t.Fatal("00E0 screen did not clear")
	} else if *u.output != [32][64]bool{} {
		t.Fatal("00E0 empty screen was not published")
	}
}

func Test_SubroutineReturn(t *testing.T) {
	c, _ := vmFixtureWithoutTick(nil)

	const newAddr = 400

	c.stack.Push(newAddr)
	c.subroutineReturn()
	if c.pc != newAddr {
		t.Fatalf("00EE subroutine return returned to incorrect address (got %d, want %d)", c.pc, newAddr)
	}
}

func Test_Jump(t *testing.T) {
	c, _ := vmFixture([]byte{0x14, 0x00})
	if c.pc != 0x400 {
		t.Fatalf("1NNN jumped to incorrect location (got %d, want %d)", c.pc, 0x400)
	}
}

func Test_SubroutineCall(t *testing.T) {
	c, _ := vmFixture([]byte{0x24, 0x00})
	if c.pc != 0x400 {
		t.Fatalf("2NNN jumped to incorrect location (got %d, want %d)", c.pc, 0x400)
	}
}

func Test_SkipEqRegConst(t *testing.T) {
	const regCont = 0x40

	// test no skip
	c, _ := vmFixtureWithoutTick([]byte{0x30, regCont})
	pcBefore := c.pc
	c.tick()
	pcAfter := c.pc

	if pcAfter - pcBefore < 2 {
		t.Fatal("3XNN incorrectly skipped")
	}

	// test skip
	c, _ = vmFixtureWithoutTick([]byte{0x30, regCont})
	c.v0 = regCont

	pcBefore = c.pc
	c.tick()
	pcAfter = c.pc

	if pcAfter - pcBefore < 4 {
		t.Fatal("3XNN failed to skip")
	}
}

func Test_SkipNotEqRegConst(t *testing.T) {
	const regCont = 0x40

	// test skip
	c, _ := vmFixtureWithoutTick([]byte{0x40, regCont})
	pcBefore := c.pc
	c.tick()
	pcAfter := c.pc

	if pcAfter - pcBefore < 4 {
		t.Fatal("4XNN failed to skip")
	}

	// test no skip
	c, _ = vmFixtureWithoutTick([]byte{0x40, regCont})
	c.v0 = regCont

	pcBefore = c.pc
	c.tick()
	pcAfter = c.pc

	if pcAfter - pcBefore < 2 {
		t.Fatal("4XNN incorrectly skipped")
	}
}

func Test_SkipEqRegReg(t *testing.T) {
	// test skip
	c, _ := vmFixtureWithoutTick([]byte{0x50, 0x10})
	c.v0 = 0x40
	c.v1 = 0x40

	pcBefore := c.pc
	c.tick()
	pcAfter := c.pc

	if pcAfter - pcBefore < 4 {
		t.Fatal("4XNN failed to skip")
	}

	// test no skip
	c, _ = vmFixtureWithoutTick([]byte{0x50, 0x10})
	c.v0 = 0x40
	c.v1 = 0x41

	pcBefore = c.pc
	c.tick()
	pcAfter = c.pc

	if pcAfter - pcBefore < 2 {
		t.Fatal("4XNN incorrectly skipped")
	}
}

func Test_SetRegisterToConstant(t *testing.T) {
	const regCont = 0x55

	c, _ := vmFixture([]byte{0x60, regCont})
	if c.v0 != regCont {
		t.Fatalf("6XNN failed to set register correctly (got %d, want %d)", c.v0, regCont)
	}
}