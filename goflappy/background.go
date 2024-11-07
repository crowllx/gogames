package main

import (
	"bytes"
	"fmt"
	"goflappy/assets"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type ParralaxBackground struct {
	fg      *ebiten.Image
	bg      *ebiten.Image
	fgSpeed float64
	bgSpeed float64
	fgPos   float64
	bgPos   float64
}

func CreateParralax() *ParralaxBackground {
	fg, _, err := image.Decode(bytes.NewReader(assets.FgImage))
	if err != nil {
		log.Fatal(err)
	}
	bg, _, err := image.Decode(bytes.NewReader(assets.BgImage))
	if err != nil {
		log.Fatal(err)
	}
	return &ParralaxBackground{
		ebiten.NewImageFromImage(fg),
		ebiten.NewImageFromImage(bg),
		0.7,
		0.5,
		0,
		0,
	}
}

func (p *ParralaxBackground) Move(x float64) {
	p.fgPos -= x * p.fgSpeed
	p.bgPos -= x * p.bgSpeed
    if p.fgPos < -float64(p.fg.Bounds().Dx()) {
        p.fgPos += float64(p.fg.Bounds().Dx()) 
    }
    if p.bgPos < -float64(p.bg.Bounds().Dx()) {
        p.bgPos += float64(p.bg.Bounds().Dx()) 
    }
}

func (p *ParralaxBackground) DrawFG(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
    fmt.Printf("pos %f\n pos2 %d", p.fgPos, p.fg.Bounds().Dx())
	opts.GeoM.Translate(p.fgPos, 0)
	screen.DrawImage(p.fg, opts)
	opts.GeoM.Translate(float64(p.fg.Bounds().Dx()), 0)
	screen.DrawImage(p.fg, opts)
}
func (p *ParralaxBackground) DrawBG(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(p.bgPos, 0)
	screen.DrawImage(p.bg, opts)
	opts.GeoM.Translate(float64(p.bg.Bounds().Dx()), 0)
	screen.DrawImage(p.bg, opts)
}
