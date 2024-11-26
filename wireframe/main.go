package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

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

var camPos = vec{0, 0, -10}

var vecs = []vec{}

var lines = []line{}

var xangle = float64(0)

var yangle = float64(0)

var zangle = float64(0)

func init() {
	loadFile("examples/fsu.edu/icosahedron.obj")
	resetCamera()
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

func resetCamera() {
	minX := math.MaxFloat64
	minY := math.MaxFloat64
	maxX := -math.MaxFloat64
	maxY := -math.MaxFloat64
	for _, vec := range vecs {
		minX = math.Min(minX, vec.x)
		minY = math.Min(minY, vec.y)
		maxX = math.Max(maxX, vec.x)
		maxY = math.Max(maxY, vec.y)
	}
	widthX := maxX - minX
	widthY := maxY - minY
	camPos = vec{minX + widthX/2, minY + widthY/2, -math.Max(widthX, widthY)}
}

func loadFile(name string) {

	bs, err := os.ReadFile(name)
	if err != nil {
		panic(err)
	}

	for _, row := range strings.Split(string(bs), "\n") {

		if strings.HasPrefix(row, "#") {
			continue
		}

		fields := strings.Fields(row)
		if len(fields) == 0 {
			continue
		}

		switch fields[0] {

		case "v":
			if len(fields) != 4 {
				panic(row)
			}

			x := parseFloat(fields[1])
			y := parseFloat(fields[2])
			z := parseFloat(fields[3])

			vecs = append(vecs, vec{x, y, z})

		case "f":

			if len(fields) != 4 {
				panic(row)
			}

			a := parseInt(fields[1])
			b := parseInt(fields[2])
			c := parseInt(fields[3])
			lines = append(lines, line{a - 1, b - 1})
			lines = append(lines, line{b - 1, c - 1})
			lines = append(lines, line{c - 1, a - 1})

		default:
			fmt.Println("ignoring:", row)

		}

	}

}

func parseFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func parseInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
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
