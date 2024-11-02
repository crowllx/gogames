package main

import (
	"github.com/crowllx/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

type Pipe struct {
	openingWidth   int
	openingStart   int
	position       int
	positionOffset int
	shapes         []*geometry.BB
}

func SpawnPipe(openingWidth int, screenWidth int, openingStart int) *Pipe {
	return &Pipe{
		openingWidth: openingWidth,
		position:     screenWidth,
		openingStart: openingStart,
		shapes: []*geometry.BB{
			&geometry.BB{
				float64(screenWidth),
				0,
				float64(screenWidth) + 30,
				float64(openingStart),
			},
			&geometry.BB{
				float64(screenWidth),
				float64(openingStart + openingWidth),
				float64(screenWidth) + 30,
				900,
			},
		},
	}
}

func (p *Pipe) Draw(screen *ebiten.Image, camera geometry.Vector) {
	vector.DrawFilledRect(
		screen,
		float32(p.position)+float32(camera.X),
		0,
		30,
		float32(p.openingStart),
		colornames.Darkcyan,
		false,
	)
	vector.DrawFilledRect(
		screen,
		float32(p.position)+float32(camera.X),
		float32(p.openingStart)+float32(p.openingWidth),
		30,
		900,
		colornames.Darkcyan,
		false,
	)

}
