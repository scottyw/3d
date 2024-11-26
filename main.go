package main

import (
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/scottyw/3d/raster"
	"golang.org/x/image/colornames"
)

const height = 1000

const width = 1000

func run() {
	cfg := opengl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width), float64(height)),
		VSync:       true,
		Undecorated: true,
	}

	// cfg.Monitor = opengl.PrimaryMonitor() // fullscreen

	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		if win.JustPressed(pixel.KeyEscape) || win.JustPressed(pixel.KeyQ) {
			return
		}
		win.Clear(colornames.Black)
		raster.Frame(width, height).Draw(win)
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
