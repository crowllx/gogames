package main

import (
	_ "fmt"

	"github.com/crowllx/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
)

func NewPlayer() Player {
	return Player{
		velocity: &geometry.Vector{X: 0, Y: 0},
		position: &geometry.Vector{X: 250, Y: 450},
		shape:    geometry.NewCircle(geometry.NewVector(250, 450), 20),
	}
}
func (p *Player) Draw(screen *ebiten.Image, camera geometry.Vector) {
	x := p.position.X + camera.X
	y := p.position.Y + camera.Y
	vector.DrawFilledCircle(
		screen,
		float32(x),
		float32(y),
		20.0,
		colornames.Coral,
		false,
	)
	vector.StrokeCircle(
		screen,
		float32(p.shape.Center().X+camera.X),
		float32(p.shape.Center().Y+camera.Y),
		float32(p.shape.Radius()),
		2,
		colornames.Darkblue,
		false,
	)
}

func (p *Player) Update() (float64, float64) {
	p.velocity.X = 100
	// apply gravity
	p.velocity.Y += float64(gravity) * dt
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.velocity.Y = -250
	}
	dx := p.velocity.X * dt
	dy := p.velocity.Y * dt
	return dx, dy
}

func (p *Player) CheckCollisions(shapes []geometry.Shape, v geometry.Vector) bool {
	p.shape.Translate(v)
	for _, s := range shapes {
		if p.shape.Collides(s) {
			p.shape.Translate(v.Neg())
			return true
		}
	}
	return false
}

func (p *Player) Move(dx, dy float64) {
	p.position.X += dx
	p.position.Y += dy
}

func (p *Player) Position() (float64, float64) {
	return p.position.X, p.position.Y
}

type Player struct {
	velocity *geometry.Vector
	position *geometry.Vector
	shape    *geometry.Circle
}
