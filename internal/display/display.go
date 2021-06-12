package display

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type Display struct {
	width, height int
	scale         int

	nextDisplay    *[32][64]bool
	currentDisplay [32][64]bool
}

func NewDisplay(width, height, scale int) *Display {
	d := &Display{
		width:  width,
		height: height,
		scale:  scale,
	}
	return d
}

func (*Display) Update() error {
	return nil
}

func (d *Display) Draw(screen *ebiten.Image) {

	if d.nextDisplay != nil {
		d.currentDisplay, d.nextDisplay = *d.nextDisplay, nil
	}

	for y := 0; y < 32; y += 1 {
		for x := 0; x < 64; x += 1 {
			var c color.Color
			if d.currentDisplay[y][x] {
				c = color.White
			} else {
				c = color.Black
			}
			screen.Set(x, y, c)
		}
	}
}

func (d *Display) Layout(outsideWidth, outsideHeight int) (int, int) {
	return d.width, d.height
}

func (d *Display) Start() error {
	ebiten.SetWindowSize(d.width*d.scale, d.height*d.scale)
	ebiten.SetWindowTitle("Hello, World!")
	return ebiten.RunGame(d)
}

func (d *Display) PublishNewDisplay(inp [32][64]bool) {
	d.nextDisplay = &inp
	fmt.Println("sent new display")
}

var inputTranslationTable = map[ebiten.Key]uint8{
	ebiten.KeyDigit1: 0x01,
	ebiten.KeyDigit2: 0x02,
	ebiten.KeyDigit3: 0x03,
	ebiten.KeyDigit4: 0x0C,
	ebiten.KeyQ:      0x04,
	ebiten.KeyW:      0x05,
	ebiten.KeyE:      0x06,
	ebiten.KeyR:      0x0D,
	ebiten.KeyA:      0x07,
	ebiten.KeyS:      0x08,
	ebiten.KeyD:      0x09,
	ebiten.KeyF:      0x0E,
	ebiten.KeyZ:      0x0A,
	ebiten.KeyX:      0x00,
	ebiten.KeyC:      0x0B,
	ebiten.KeyV:      0x0F,
}

func (d *Display) GetPressedKeys() []uint8 {
	var o []uint8

	for _, p := range inpututil.PressedKeys() {
		if x, ok := inputTranslationTable[p]; ok {
			o = append(o, x)
		}
	}

	return o
}
