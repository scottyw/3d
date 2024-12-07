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
		if win.JustPressed(pixel.KeyEscape) {
			return
		}

		switch {
		case win.Pressed(pixel.KeyW):
			raster.Forward()
		case win.Pressed(pixel.KeyS):
			raster.Back()
		case win.Pressed(pixel.KeyUp):
			raster.Forward()
		case win.Pressed(pixel.KeyDown):
			raster.Back()
		case win.Pressed(pixel.KeyD):
			raster.StrafeRight()
		case win.Pressed(pixel.KeyA):
			raster.StrafeLeft()
		case win.Pressed(pixel.KeyRight):
			raster.TurnRight()
		case win.Pressed(pixel.KeyLeft):
			raster.TurnLeft()
		case win.Pressed(pixel.KeyQ):
			raster.Up()
		case win.Pressed(pixel.KeyE):
			raster.Down()

		}

		win.Clear(colornames.Black)
		raster.Frame(width, height).Draw(win)
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
