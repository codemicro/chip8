package display

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Display struct {
	width, height int
	scale int

	nextDisplay *[32][64]bool
	currentDisplay [32][64]bool
}

func NewDisplay(width, height, scale int) *Display {
	d := &Display{
		width: width,
		height: height,
		scale: scale,
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
	ebiten.SetWindowSize(d.width * d.scale, d.height * d.scale)
	ebiten.SetWindowTitle("Hello, World!")
	return ebiten.RunGame(d)
}

func (d *Display) PublishNewDisplay(inp [32][64]bool) {
	d.nextDisplay = &inp
	fmt.Println("sent new display")
}

func (d *Display) GetPressedKeys() []uint8 {
	return nil
}