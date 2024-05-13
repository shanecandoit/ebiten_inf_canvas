package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 600
	screenHeight = 400
)

var (
	bgImage *ebiten.Image
)

func init() {
	// Decode an image from the image file's byte slice.
	// image_name := "5-mm-dot-paper-a4.png"
	image_name := "5-mm-white-dot-black-paper-a4.png"

	img, err := os.ReadFile(image_name)
	if err != nil {
		log.Fatal(err)
	}
	image_decoded, _, err := image.Decode(bytes.NewReader(img))
	if err != nil {
		log.Fatal(err)
	}
	bgImage = ebiten.NewImageFromImage(image_decoded)
}

type viewport struct {
	x16 int
	y16 int
}

func (p *viewport) Move() {
	s := bgImage.Bounds().Size()
	maxX16 := s.X * 16
	maxY16 := s.Y * 16

	// p.x16 += s.X / 32
	// p.y16 += s.Y / 32
	p.x16 %= maxX16
	p.y16 %= maxY16
}

func (p *viewport) Position() (int, int) {
	return p.x16, p.y16
}

type Game struct {
	viewport viewport
}

func (g *Game) Update() error {
	g.viewport.Move()

	// mouse wheel vertical scroll
	_, vert := ebiten.Wheel()
	if vert != 0 {
		g.viewport.y16 += int(vert) * 16 * 10
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	x16, y16 := g.viewport.Position()
	offsetX, offsetY := float64(-x16)/16, float64(-y16)/16

	// Draw bgImage on the screen repeatedly.
	const repeat = 3
	w, h := bgImage.Bounds().Dx(), bgImage.Bounds().Dy()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), float64(h*j))
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(bgImage, op)
		}
	}

	text := fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS())
	text += fmt.Sprintf("\npos: %0.2f %0.2f", offsetX, offsetY)
	ebitenutil.DebugPrint(screen, text)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {

	// ebiten.SetWindowSize(screenWidth, screenHeight)
	return outsideWidth, outsideHeight
	// return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Infinite Scroll (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
