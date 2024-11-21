package main

import (
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

var height = 1000
var width = 1000

type point struct {
	x float64
	y float64
	z float64
}

type line struct {
	p1 int
	p2 int
}

var points = []point{}

var lines = []line{}

var xangle = float64(0)

var yangle = float64(0)

var zangle = float64(0)

func init() {
	addCuboid(-100, -100, -100, 100, 100, 100)

	addCuboid(200, -100, -100, 300, 100, 100)
	addCuboid(-100, 200, -100, 100, 300, 100)
	addCuboid(-100, -100, 200, 100, 100, 300)

	addCuboid(-200, -100, -100, -300, 100, 100)
	addCuboid(-100, -200, -100, 100, -300, 100)
	addCuboid(-100, -100, -200, 100, 100, -300)

}

func addCuboid(x1, y1, z1, x2, y2, z2 float64) {
	i := len(points)
	points = append(points, point{x1, y1, z1})
	points = append(points, point{x2, y1, z1})
	points = append(points, point{x2, y2, z1})
	points = append(points, point{x1, y2, z1})
	points = append(points, point{x1, y1, z2})
	points = append(points, point{x2, y1, z2})
	points = append(points, point{x2, y2, z2})
	points = append(points, point{x1, y2, z2})
	lines = append(lines, line{i, i + 1})
	lines = append(lines, line{i + 1, i + 2})
	lines = append(lines, line{i + 2, i + 3})
	lines = append(lines, line{i + 3, i + 0})
	lines = append(lines, line{i + 4, i + 5})
	lines = append(lines, line{i + 5, i + 6})
	lines = append(lines, line{i + 6, i + 7})
	lines = append(lines, line{i + 7, i + 4})
	lines = append(lines, line{i + 0, i + 4})
	lines = append(lines, line{i + 1, i + 5})
	lines = append(lines, line{i + 2, i + 6})
	lines = append(lines, line{i + 3, i + 7})
}

func frame() *imdraw.IMDraw {

	imd := imdraw.New(nil)
	imd.Color = colornames.Lawngreen

	updatedPoints := []point{}
	for _, point := range points {
		point = rotateX(point, xangle)
		point = rotateY(point, yangle)
		point = rotateZ(point, zangle)
		point = perspective(point)
		point = center(point)
		updatedPoints = append(updatedPoints, point)
	}

	for _, line := range lines {
		p1 := updatedPoints[line.p1]
		p2 := updatedPoints[line.p2]

		imd.Push(pixel.V(p1.x, p1.y))
		imd.Push(pixel.V(p2.x, p2.y))
		imd.Line(2)

	}

	xangle += 0.003
	yangle += 0.005
	zangle += 0.007

	return imd
}

func center(p point) point {
	return point{
		p.x + float64(width)/2,
		p.y + float64(height)/2,
		p.z,
	}
}

func perspective(p point) point {
	distance := float64(width) / 2
	zdelta := distance / (distance + p.z)
	return point{
		p.x * zdelta,
		p.y * zdelta,
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
		frame().Draw(win)
		win.Update()
	}
}

func main() {
	opengl.Run(run)
}
