package raster

import (
	"fmt"
	"image/color"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/imdraw"
)

type vec2 struct {
	x float64
	y float64
}

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

var cameraPosition vec

// Camera always points directly into the Z axis
var cameraDirection = vec{0, 0, 1}

// Light source is infinitely far away and points in this direction
var lightSource = normalize(vec{1, -1, 1})

var vecs = []vec{}

var triangles = []triangle{}

var xangle = float64(0)

var yangle = float64(0)

var zangle = float64(0)

func init() {
	loadFile("examples/fsu.edu/icosahedron.obj")
	resetCameraPosition()
}

func Frame(width, height int) *imdraw.IMDraw {

	imd := imdraw.New(nil)

	// Rotate and translate each vertex of the loaded scene in world space
	updated := []vec{}
	for _, vec := range vecs {
		vec = rotateX(vec, xangle)
		vec = rotateY(vec, yangle)
		vec = rotateZ(vec, zangle)
		vec = translateRelativeToCamera(vec)
		updated = append(updated, vec)
	}

	// Sort triangles by Z position
	slices.SortStableFunc(triangles, func(a, b triangle) int {
		centreA := centre(a, updated)
		centreB := centre(b, updated)
		diff := centreB.z - centreA.z
		if diff == 0 {
			return 0
		}
		if diff < 0 {
			return -1
		}
		return 1
	})

	// Render each triangle of the loaded scene
	for _, triangle := range triangles {

		// Project this triangle's "a" vertex onto the 2D plane
		a := projectInto2D(updated[triangle.a])
		ax := (a.x + 1) * float64(width) / 2
		ay := (a.y + 1) * float64(height) / 2

		// Project this triangle's "b" vertex onto the 2D plane
		b := projectInto2D(updated[triangle.b])
		bx := (b.x + 1) * float64(width) / 2
		by := (b.y + 1) * float64(height) / 2

		// Project this triangle's "c" vertex onto the 2D plane
		c := projectInto2D(updated[triangle.c])
		cx := (c.x + 1) * float64(width) / 2
		cy := (c.y + 1) * float64(height) / 2

		// Draw the triangle onscreen

		normal := findNormal(triangle, updated)
		light := dot(lightSource, normal)
		brightness := uint8(0x80 - (light * 0x60))
		imd.Color = color.RGBA{brightness, brightness, brightness, 0xff}

		imd.Push(pixel.V(ax, ay))
		imd.Push(pixel.V(bx, by))
		imd.Push(pixel.V(cx, cy))
		imd.Polygon(0)

	}

	// Rotate the loaded scene
	yangle += 0.005

	return imd
}

func projectInto2D(p vec) vec2 {
	zoffset := focal / (focal + p.z)
	return vec2{
		p.x * zoffset,
		p.y * zoffset,
	}
}

func translateRelativeToCamera(p vec) vec {
	return vec{
		p.x - cameraPosition.x,
		p.y - cameraPosition.y,
		p.z - cameraPosition.z,
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
		p.x*cos - p.y*sin,
		p.x*sin + p.y*cos,
		p.z,
	}
}

func resetCameraPosition() {
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
	cameraPosition = vec{minX + (widthX * 0.5), minY + (widthY * 0.75), -1.2 * math.Max(widthX, widthY)}
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

			if len(fields) < 4 {
				panic(row)
			}

			for x := 2; x < len(fields)-1; x++ {
				a := parseInt(fields[1])
				b := parseInt(fields[x])
				c := parseInt(fields[x+1])

				t := triangle{a - 1, b - 1, c - 1}
				triangles = append(triangles, t)
			}

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

func findNormal(t triangle, vecs []vec) vec {
	ba := sub(vecs[t.b], vecs[t.a])
	ca := sub(vecs[t.c], vecs[t.a])
	cross := cross(ba, ca)
	return normalize(cross)
}

func centre(t triangle, vecs []vec) vec {
	v1 := vecs[t.a]
	v2 := vecs[t.b]
	v3 := vecs[t.c]
	return vec{
		(v1.x + v2.x + v3.x) / 3,
		(v1.y + v2.y + v3.y) / 3,
		(v1.z + v2.z + v3.z) / 3,
	}
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
