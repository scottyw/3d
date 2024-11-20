package main

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
)

var height = 1000
var width = 1000

func main() {
	opengl.Run(run)
}

func run() {
	cfg := opengl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width), float64(height)),
		VSync:       true,
		Undecorated: true,
	}

	// fullscreen
	// cfg.Monitor = opengl.PrimaryMonitor()

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	c := win.Bounds().Center()

	for !win.Closed() {
		if win.JustPressed(pixel.KeyEscape) || win.JustPressed(pixel.KeyQ) {
			return
		}

		win.Clear(color.White)

		p := pixel.PictureDataFromImage(frame())

		pixel.NewSprite(p, p.Bounds()).Draw(win, pixel.IM.Moved(c))

		win.Update()
	}
}

func frame() image.Image {
	// m := image.NewRGBA(image.Rect(0, 0, width, height))

	dc := gg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(0, 0, 0)
	dc.Fill()

	return dc.Image()
}
