package main

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
)

var height = float64(1000)
var width = float64(1000)

type point struct {
	x float64
	y float64
	z float64
}

type line struct {
	p1 int
	p2 int
}

var points = []point{
	{1, 1, 1},
	{1, -1, 1},
	{-1, -1, 1},
	{-1, 1, 1},
	{1, 1, -1},
	{1, -1, -1},
	{-1, -1, -1},
	{-1, 1, -1},
}

var lines = []line{
	{0, 1},
	{1, 2},
	{2, 3},
	{3, 0},
	{4, 5},
	{5, 6},
	{6, 7},
	{7, 4},
	{0, 4},
	{1, 5},
	{2, 6},
	{3, 7},
}

var xangle = float64(0)

var yangle = float64(0)

var zangle = float64(0)

func frame() image.Image {
	dc := gg.NewContext(int(width), int(height))
	dc.SetColor(colornames.Lawngreen)
	dc.SetLineWidth(2)

	updatedPoints := []point{}
	for _, point := range points {
		point = rotateX(point, xangle)
		point = rotateY(point, yangle)
		point = rotateZ(point, zangle)
		point = perspective(point)
		point = scale(point, 200)
		point = move(point, width, height)
		updatedPoints = append(updatedPoints, point)
	}

	for _, line := range lines {
		p1 := updatedPoints[line.p1]
		p2 := updatedPoints[line.p2]
		dc.DrawLine(p1.x, p1.y, p2.x, p2.y)
	}

	xangle += 0.003
	yangle += 0.005
	zangle += 0.007

	dc.Stroke()
	return dc.Image()
}

func scale(p point, s float64) point {
	return point{
		p.x * s,
		p.y * s,
		p.z * s,
	}
}

func move(p point, w, h float64) point {
	return point{
		p.x + w/2,
		p.y + h/2,
		p.z,
	}
}

func perspective(p point) point {
	if p.z >= 5 {
		return point{
			p.x,
			p.y,
			0,
		}
	}
	z := (3 / (5 - p.z))
	return point{
		p.x * z,
		p.y * z,
		0,
	}
}

func rotateX(p point, theta float64) point {
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	return point{
		p.x,
		p.y*cos - p.z*sin,
		p.y*sin + p.z*cos,
	}
}

func rotateY(p point, theta float64) point {
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	return point{
		p.x*cos - p.z*sin,
		p.y,
		p.x*sin + p.z*cos,
	}
}

func rotateZ(p point, theta float64) point {
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	return point{
		p.x,
		p.y*cos - p.z*sin,
		p.y*sin + p.z*cos,
	}
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

		win.Clear(color.Black)

		p := pixel.PictureDataFromImage(frame())

		pixel.NewSprite(p, p.Bounds()).Draw(win, pixel.IM.Moved(c))

		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
