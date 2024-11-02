package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"time"

	"github.com/crowllx/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/gomono"
)

// reset player x value to avoid an infinitely increasing x
// func (g *Game) resetXPosition() {
// 	g.player.position.X += g.camera.X
// 	for _, pipe := range g.pipes {
// 		pipe.position += int(g.camera.X)
// 	}
// 	g.camera.X = -g.player.position.X + float64(g.width)/4
// }

func randRange(min, max int) int {
	return rand.IntN(max+1-min) + min
}
func NewGame() *Game {
	g := Game{
		width:     1600,
		height:    900,
		gravity:   800,
		player:    NewPlayer(),
		camera:    geometry.Vector{X: 0, Y: 0},
		pipes:     []*Pipe{},
		colliders: make([]geometry.Shape, 0),
		canSpawn:  true,
		startTime: time.Now(),
		duration:  time.Duration(0),
		lastSpawn: 0,
	}
	return &g
}

type Game struct {
	width      int
	height     int
	gravity    int
	player     Player
	camera     geometry.Vector
	pipes      []*Pipe
	colliders  []geometry.Shape
	canSpawn   bool
	startTime  time.Time
	duration   time.Duration
	textSource *text.GoTextFaceSource
	lastSpawn  int
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"x: %f\n y: %f\n # pipes: %d\n # colliders: %d",
			g.player.position.X,
			g.player.position.Y,
			len(g.pipes),
			len(g.colliders),
		),
	)
	// object drawing
	g.player.Draw(screen, g.camera)
	for _, pipe := range g.pipes {
		pipe.Draw(screen, g.camera)
	}

	// text display
	face := &text.GoTextFace{
		Source: g.textSource,
		Size:   24,
	}
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(50, 50)
	minutes := g.duration.Seconds() / 60
	seconds := math.Mod(g.duration.Seconds(), 60)
	str := fmt.Sprintf("%02d:%02d", int(minutes), int(seconds))
	text.Draw(screen, str, face, opts)

}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return g.width, g.height
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.duration = time.Now().Sub(g.startTime)
	// if g.player.position.X > 1600 {
	// 	g.resetXPosition()
	// }
	g.camera.X = -g.player.position.X + float64(g.width/4.0)
	dx, dy := g.player.Update()

	if g.player.CheckCollisions(g.colliders, geometry.NewVector(dx, 0)) {
		dx = 0
	}

	if g.player.CheckCollisions(g.colliders, geometry.NewVector(0, dy)) {
		dy = 0
	}

	g.player.Move(dx, dy)

	fmt.Printf("%d\n", g.lastSpawn-int(1600-g.camera.X))
	if g.canSpawn && g.lastSpawn-int(1600-g.camera.X) < -250 {
		newPipe := SpawnPipe(randRange(150, 300), int(1600-g.camera.X), randRange(100, 500))
		g.lastSpawn = int(1600 - g.camera.X)
		g.pipes = append(g.pipes, newPipe)
		for _, shape := range newPipe.shapes {
			g.colliders = append(g.colliders, shape)
		}
		g.canSpawn = false
		duration := time.Duration(rand.IntN(3000) + 3000)
		timer := time.NewTimer(duration * time.Millisecond)
		fmt.Println(duration)
		go func() {
			<-timer.C
			g.canSpawn = true
		}()

		log.Printf(
			"Pipe spawned, width: %d, start: %d, next spawn: %s",
			newPipe.openingWidth,
			newPipe.openingStart,
			duration,
		)
	}

    //remove off screen pipes
	if g.pipes[0].position+30+int(g.camera.X) < 0 {
		g.pipes = g.pipes[1:]
		g.colliders = g.colliders[2:]
	}
	// camera
	return nil
}

var _ ebiten.Game = &Game{}

func main() {
	g := NewGame()
	textSource, err := text.NewGoTextFaceSource(bytes.NewReader(gomono.TTF))
	if err != nil {
		panic(err)
	}
	g.textSource = textSource
	log.SetFlags(log.Ltime)
	ebiten.SetWindowSize(g.width, g.height)
	ebiten.SetWindowTitle("GoFlappy")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
