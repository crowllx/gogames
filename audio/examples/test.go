package main

import (
	_ "embed"
	"fmt"
	"log"
	"math"

	"github.com/crowllx/gogames/audio"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"
)

var (
	//go:embed "Sketchbook 2023-11-29.ogg"
	bgm_ogg []byte
)

type object struct {
	x, y float64
}

func (o *object) Position() (float64, float64) {
	return o.x, o.y
}

type Game struct {
	bgm          *audio.AudioSource
	count        int
	xpos         float64
	screenWidth  int
	screenHeight int
	panning      float64
}

var ebitenImage *ebiten.Image

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, fmt.Sprintf("panning: %.2f", g.panning))
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(g.xpos-float64(ebitenImage.Bounds().Dx()/2), float64(g.screenWidth)/2)
	screen.DrawImage(ebitenImage, opts)
}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.screenWidth, g.screenHeight
}
func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.count++
	r := float64(g.count) * ((1.0 / 60.0) * 2 * math.Pi) * 0.1
	g.xpos = (float64(g.screenWidth) / 2) + math.Cos(r)*(float64(g.screenWidth)/2)
	g.panning = lerp(-1, 1, g.xpos/float64(g.screenWidth))
	g.bgm.SetPan("bgm", g.panning)

	return nil
}

func NewGame() *Game {
	ebitenImage = ebiten.NewImage(20, 20)
	ebitenImage.Fill(colornames.Aliceblue)
	obj := &object{0, 0}

	src := audio.NewSource(obj)
	err := src.AddController("bgm", bgm_ogg, false)
	if err != nil {
		panic(err)
	}
	src.Play("bgm")
	return &Game{
		bgm:          src,
		screenWidth:  640,
		screenHeight: 480,
        panning: 0,
	}
}

var _ ebiten.Game = &Game{}

func main() {
	g := NewGame()

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
