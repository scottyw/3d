package raster

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
	"golang.org/x/image/colornames"
)

type vec struct {
	x float64
	y float64
	z float64
}

type triangle struct {
	a int
	b int
	c int
}

const focal = float64(1)

var camPos vec

var vecs = []vec{}

var triangles = []triangle{}

var normals = []vec{}

var xangle = float64(0)

var yangle = float64(0)

var zangle = float64(0)

func init() {
	loadFile("examples/fsu.edu/icosahedron.obj")
	resetCamera()
}

func Frame(width, height int) *imdraw.IMDraw {

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

	for _, triangle := range triangles {
		a := updated[triangle.a]
		ax := (a.x + 1) * float64(width) / 2
		ay := (a.y + 1) * float64(height) / 2
		b := updated[triangle.b]
		bx := (b.x + 1) * float64(width) / 2
		by := (b.y + 1) * float64(height) / 2
		c := updated[triangle.c]
		cx := (c.x + 1) * float64(width) / 2
		cy := (c.y + 1) * float64(height) / 2
		// if normals[i].z >= 0 {
		imd.Push(pixel.V(ax, ay))
		imd.Push(pixel.V(bx, by))
		imd.Push(pixel.V(cx, cy))
		imd.Polygon(2)
		// }
	}

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

			t := triangle{a - 1, b - 1, c - 1}
			triangles = append(triangles, t)
			normals = append(normals, findNormal(t))

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

func findNormal(t triangle) vec {
	ba := sub(vecs[t.b], vecs[t.a])
	ca := sub(vecs[t.c], vecs[t.a])
	cross := cross(ba, ca)
	return normalize(cross)
}

func sub(a, b vec) vec {
	return vec{a.x - b.x, a.y - b.y, a.z - b.z}
}

func dot(a, b vec) float64 {
	return a.x*b.x + a.y*b.y + a.z*b.z
}

func cross(a, b vec) vec {
	return vec{
		a.y*b.z - a.z*b.y,
		a.z*b.x - a.x*b.z,
		a.x*b.y - a.y*b.x,
	}
}

func length(a vec) float64 {
	return math.Sqrt(dot(a, a))
}

func normalize(a vec) vec {
	l := length(a)
	return vec{a.x / l, a.y / l, a.z / l}
}

func Right() {
	yangle += 0.01
}

func Left() {
	yangle -= 0.01
}

func Up() {
	xangle += 0.01
}

func Down() {
	xangle -= 0.01

}
