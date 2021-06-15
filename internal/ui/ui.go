package ui

import (
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"strconv"
	"strings"
)

type UI struct {
	Debug bool

	audioPlayer *audio.Player
	toneFrequency int

	fgColour color.Color
	bgColour color.Color

	width, height int
	scale         int
	windowTitle   string

	nextDisplay    *[32][64]bool
	currentDisplay [32][64]bool
}

func NewUI(width, height, scale int, windowTitle string, toneFrequency int, fgColour, bgColour string) (*UI, error) {

	fg, err := hexStringToColor(fgColour)
	if err != nil {
		return nil, err
	}

	bg, err := hexStringToColor(bgColour)
	if err != nil {
		return nil, err
	}

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

		fgColour: fg,
		bgColour: bg,
	}
	return d, nil
}

func hexStringToColor(hx string) (color.Color, error) {
	hx = strings.TrimPrefix(hx, "#")

	if !(len(hx) == 6 || len(hx) == 3) {
		return nil, errors.New("invalid colour")
	}

	ston := func(y string) (uint8, error) {
		l, err := strconv.ParseInt(y, 16, 9) // 9 seems to play nicely with uint8, whereas 8 does not
		return uint8(l), err
	}

	var rs, gs, bs string

	switch len(hx) {
	case 6:
		rs = hx[0:2]
		gs = hx[2:4]
		bs = hx[4:6]
	case 3:
		dbf := func(x uint8) string {
			y := string(x)
			return y + y
		}

		rs = dbf(hx[0])
		gs = dbf(hx[1])
		bs = dbf(hx[2])
	default:
		return nil, errors.New("invalid colour")
	}

	r, err := ston(rs)
	if err != nil {
		return nil, err
	}

	g, err := ston(gs)
	if err != nil {
		return nil, err
	}

	b, err := ston(bs)
	if err != nil {
		return nil, err
	}

	return color.RGBA{
		R: r,
		G: g,
		B: b,
	}, nil
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
					c = d.fgColour
				}
			} else {
				c = d.bgColour
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