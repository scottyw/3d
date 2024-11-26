package main

import (
	"math"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type vec struct {
	x float64
	y float64
	z float64
}

type line struct {
	p1 int
	p2 int
}

const height = 1000

const width = 1000

const focal = float64(1)

var camPos = vec{-50, 50, -1000}

var vecs = []vec{}

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
	i := len(vecs)
	vecs = append(vecs, vec{x1, y1, z1})
	vecs = append(vecs, vec{x2, y1, z1})
	vecs = append(vecs, vec{x2, y2, z1})
	vecs = append(vecs, vec{x1, y2, z1})
	vecs = append(vecs, vec{x1, y1, z2})
	vecs = append(vecs, vec{x2, y1, z2})
	vecs = append(vecs, vec{x2, y2, z2})
	vecs = append(vecs, vec{x1, y2, z2})
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

	updated := []vec{}
	for _, vec := range vecs {
		vec = rotateX(vec, xangle)
		vec = rotateY(vec, yangle)
		vec = rotateZ(vec, zangle)
		vec = translate(vec)
		vec = projection(vec)
		updated = append(updated, vec)
	}

	for _, line := range lines {
		p1 := updated[line.p1]
		p1x := (p1.x + 1) * float64(width) / 2
		p1y := (p1.y + 1) * float64(height) / 2
		p2 := updated[line.p2]
		p2x := (p2.x + 1) * float64(width) / 2
		p2y := (p2.y + 1) * float64(height) / 2
		imd.Push(pixel.V(p1x, p1y))
		imd.Push(pixel.V(p2x, p2y))
		imd.Line(2)
	}

	xangle += 0.003
	yangle += 0.005
	zangle += 0.007

	return imd
}

func translate(p vec) vec {
	return vec{
		p.x - camPos.x,
		p.y - camPos.y,
		p.z - camPos.z,
	}
}

func projection(p vec) vec {
	zoffset := focal / (focal + p.z)
	return vec{
		p.x * zoffset,
		p.y * zoffset,
		0,
	}
}

func rotateX(p vec, theta float64) vec {
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	return vec{
		p.x,
		p.y*cos - p.z*sin,
		p.y*sin + p.z*cos,
	}
}

func rotateY(p vec, theta float64) vec {
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	return vec{
		p.x*cos - p.z*sin,
		p.y,
		p.x*sin + p.z*cos,
	}
}

func rotateZ(p vec, theta float64) vec {
	sin := math.Sin(theta)
	cos := math.Cos(theta)
	return vec{
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
