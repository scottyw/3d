package main

import (
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/scottyw/3d/raster"
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
		if win.JustPressed(pixel.KeyEscape) {
			return
		}
		win.Clear(color.RGBA{0x20, 0x20, 0x20, 0xff})
		raster.Frame(width, height).Draw(win)
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
