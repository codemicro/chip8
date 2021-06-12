package vm

import (
	"encoding/binary"
	"fmt"
)

// getAddressFromCIR gets the second, third and fourth nibbles (0NNN) from the current instruction and returns it as a
// uint16
func (c *Chip8) getAddressFromCIR() uint16 {
	return binary.BigEndian.Uint16([]byte{
		c.cir[0] & 0x0F,
		c.cir[1],
	})
}

func (c *Chip8) get8bitConstant() byte {
	return c.cir[1]
}

func (c *Chip8) get4BitConstant() byte {
	return c.cir[1] & 0x0F
}

func (c *Chip8) getRegisterPointer(register byte) *byte {
	switch register {
	case 0x00:
		return &c.v0
	case 0x01:
		return &c.v1
	case 0x02:
		return &c.v2
	case 0x03:
		return &c.v3
	case 0x04:
		return &c.v4
	case 0x05:
		return &c.v5
	case 0x06:
		return &c.v6
	case 0x07:
		return &c.v7
	case 0x08:
		return &c.v8
	case 0x09:
		return &c.v9
	case 0x0A:
		return &c.va
	case 0x0B:
		return &c.vb
	case 0x0C:
		return &c.vc
	case 0x0D:
		return &c.vd
	case 0x0E:
		return &c.ve
	case 0x0F:
		return &c.vf
	default:
		panic(fmt.Errorf("unknown register 0x%x", register))
	}
}