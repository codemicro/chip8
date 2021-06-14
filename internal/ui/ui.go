package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
)

type UI struct {
	Debug bool

	audioPlayer *audio.Player
	toneFrequency int

	width, height int
	scale         int
	windowTitle   string

	nextDisplay    *[32][64]bool
	currentDisplay [32][64]bool
}

func NewUI(width, height, scale int, windowTitle string, toneFrequency int) (*UI, error) {

	p, err := audio.NewPlayer(audioContext, &stream{toneFrequency: toneFrequency})
	if err != nil {
		return nil, err
	}

	d := &UI{
		width:  width,
		height: height,
		scale:  scale,
		windowTitle: windowTitle,

		audioPlayer: p,
		toneFrequency: toneFrequency,
	}
	return d, nil
}

func (*UI) Update() error {
	return nil
}

func (d *UI) Draw(screen *ebiten.Image) {

	if d.nextDisplay != nil {
		d.currentDisplay, d.nextDisplay = *d.nextDisplay, nil
	}

	for y := 0; y < 32; y += 1 {
		for x := 0; x < 64; x += 1 {
			var c color.Color
			if d.currentDisplay[y][x] {
				if x % 2 == 0 && d.Debug {
					c = color.RGBA{
						R: 255,
					}
				} else {
					// c = color.White
					c = color.RGBA{
						R: 0x3D,
						G: 0x80,
						B: 0x26,
					}
				}
			} else {
				// c = color.Black
				c = color.RGBA{
					R: 0xF9,
					G: 0xFF,
					B: 0xB3,
				}
			}
			screen.Set(x, y, c)
		}
	}
}

func (d *UI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return d.width, d.height
}

func (d *UI) Start() error {
	ebiten.SetWindowSize(d.width*d.scale, d.height*d.scale)
	ebiten.SetWindowTitle(d.windowTitle)
	return ebiten.RunGame(d)
}

func (d *UI) PublishNewDisplay(inp [32][64]bool) {
	x := inp // copy
	d.nextDisplay = &x
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

func (d *UI) GetPressedKeys() []uint8 {
	var o []uint8

	for _, p := range inpututil.PressedKeys() {
		if x, ok := inputTranslationTable[p]; ok {
			o = append(o, x)
		}
	}

	return o
}

func (d *UI) StartTone() {
	if !d.audioPlayer.IsPlaying() {
		d.audioPlayer.Play()
	}
}

func (d *UI) StopTone() {
	if d.audioPlayer.IsPlaying() {
		d.audioPlayer.Pause()
	}
}