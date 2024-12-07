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

type vec4 struct {
	x float64
	y float64
	z float64
	w float64
}

type m4 struct {
	m00, m01, m02, m03,
	m10, m11, m12, m13,
	m20, m21, m22, m23,
	m30, m31, m32, m33 float64
}

type triangle struct {
	a int
	b int
	c int
}

const focal = float64(1)

var camPos vec

var camDir vec

var vecs = []vec{}

var triangles = []triangle{}

var normals = []vec{}

func init() {
	loadFile("examples/axis.obj")
	resetCamera()
}

func Frame(width, height int) *imdraw.IMDraw {

	right := normalize(cross(vec{0, 1, 0}, camDir))

	up := cross(camDir, right)

	lookAt := m4{
		right.x, right.y, right.z, 0,
		up.x, up.y, up.z, 0,
		camDir.x, camDir.y, camDir.z, 0,
		0, 0, 0, 1,
	}

	lookAt = multiply(lookAt, m4{
		1, 0, 0, -camPos.x,
		0, 1, 0, -camPos.y,
		0, 0, 1, -camPos.z,
		0, 0, 0, 1,
	})

	updated := []vec4{}
	for _, v := range vecs {
		v4 := vec4{v.x, v.y, v.z, 1}
		v4 = matrix(lookAt, v4)
		v4 = projection(v4)
		updated = append(updated, v4)
	}

	imd := imdraw.New(nil)
	imd.Color = colornames.Lawngreen

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

func projection(p vec4) vec4 {
	zoffset := focal / (focal + p.z)
	return vec4{
		x: p.x * zoffset,
		y: p.y * zoffset,
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
	camDir = normalize(sub(vec{0, 0, 0}, camPos))
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

func add(a, b vec) vec {
	return vec{a.x + b.x, a.y + b.y, a.z + b.z}
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

func matrix(m m4, v vec4) vec4 {
	return vec4{
		m.m00*v.x + m.m01*v.y + m.m02*v.z + m.m03*v.w,
		m.m10*v.x + m.m11*v.y + m.m12*v.z + m.m13*v.w,
		m.m20*v.x + m.m21*v.y + m.m22*v.z + m.m23*v.w,
		m.m30*v.x + m.m31*v.y + m.m32*v.z + m.m33*v.w,
	}
}

func multiply(m, m2 m4) m4 {
	return m4{
		m00: m.m00*m2.m00 + m.m01*m2.m10 + m.m02*m2.m20 + m.m03*m2.m30,
		m01: m.m00*m2.m01 + m.m01*m2.m11 + m.m02*m2.m21 + m.m03*m2.m31,
		m02: m.m00*m2.m02 + m.m01*m2.m12 + m.m02*m2.m22 + m.m03*m2.m32,
		m03: m.m00*m2.m03 + m.m01*m2.m13 + m.m02*m2.m23 + m.m03*m2.m33,
		m10: m.m10*m2.m00 + m.m11*m2.m10 + m.m12*m2.m20 + m.m13*m2.m30,
		m11: m.m10*m2.m01 + m.m11*m2.m11 + m.m12*m2.m21 + m.m13*m2.m31,
		m12: m.m10*m2.m02 + m.m11*m2.m12 + m.m12*m2.m22 + m.m13*m2.m32,
		m13: m.m10*m2.m03 + m.m11*m2.m13 + m.m12*m2.m23 + m.m13*m2.m33,
		m20: m.m20*m2.m00 + m.m21*m2.m10 + m.m22*m2.m20 + m.m23*m2.m30,
		m21: m.m20*m2.m01 + m.m21*m2.m11 + m.m22*m2.m21 + m.m23*m2.m31,
		m22: m.m20*m2.m02 + m.m21*m2.m12 + m.m22*m2.m22 + m.m23*m2.m32,
		m23: m.m20*m2.m03 + m.m21*m2.m13 + m.m22*m2.m23 + m.m23*m2.m33,
		m30: m.m30*m2.m00 + m.m31*m2.m10 + m.m32*m2.m20 + m.m33*m2.m30,
		m31: m.m30*m2.m01 + m.m31*m2.m11 + m.m32*m2.m21 + m.m33*m2.m31,
		m32: m.m30*m2.m02 + m.m31*m2.m12 + m.m32*m2.m22 + m.m33*m2.m32,
		m33: m.m30*m2.m03 + m.m31*m2.m13 + m.m32*m2.m23 + m.m33*m2.m33,
	}
}

// From https://stackoverflow.com/questions/1148309/inverting-a-4x4-matrix
func invert(m m4) m4 {

	var A2323 = m.m22*m.m33 - m.m23*m.m32
	var A1323 = m.m21*m.m33 - m.m23*m.m31
	var A1223 = m.m21*m.m32 - m.m22*m.m31
	var A0323 = m.m20*m.m33 - m.m23*m.m30
	var A0223 = m.m20*m.m32 - m.m22*m.m30
	var A0123 = m.m20*m.m31 - m.m21*m.m30
	var A2313 = m.m12*m.m33 - m.m13*m.m32
	var A1313 = m.m11*m.m33 - m.m13*m.m31
	var A1213 = m.m11*m.m32 - m.m12*m.m31
	var A2312 = m.m12*m.m23 - m.m13*m.m22
	var A1312 = m.m11*m.m23 - m.m13*m.m21
	var A1212 = m.m11*m.m22 - m.m12*m.m21
	var A0313 = m.m10*m.m33 - m.m13*m.m30
	var A0213 = m.m10*m.m32 - m.m12*m.m30
	var A0312 = m.m10*m.m23 - m.m13*m.m20
	var A0212 = m.m10*m.m22 - m.m12*m.m20
	var A0113 = m.m10*m.m31 - m.m11*m.m30
	var A0112 = m.m10*m.m21 - m.m11*m.m20

	var det = m.m00*(m.m11*A2323-m.m12*A1323+m.m13*A1223) -
		m.m01*(m.m10*A2323-m.m12*A0323+m.m13*A0223) +
		m.m02*(m.m10*A1323-m.m11*A0323+m.m13*A0123) -
		m.m03*(m.m10*A1223-m.m11*A0223+m.m12*A0123)

	det = 1 / det

	return m4{
		m00: det * (m.m11*A2323 - m.m12*A1323 + m.m13*A1223),
		m01: det * -(m.m01*A2323 - m.m02*A1323 + m.m03*A1223),
		m02: det * (m.m01*A2313 - m.m02*A1313 + m.m03*A1213),
		m03: det * -(m.m01*A2312 - m.m02*A1312 + m.m03*A1212),
		m10: det * -(m.m10*A2323 - m.m12*A0323 + m.m13*A0223),
		m11: det * (m.m00*A2323 - m.m02*A0323 + m.m03*A0223),
		m12: det * -(m.m00*A2313 - m.m02*A0313 + m.m03*A0213),
		m13: det * (m.m00*A2312 - m.m02*A0312 + m.m03*A0212),
		m20: det * (m.m10*A1323 - m.m11*A0323 + m.m13*A0123),
		m21: det * -(m.m00*A1323 - m.m01*A0323 + m.m03*A0123),
		m22: det * (m.m00*A1313 - m.m01*A0313 + m.m03*A0113),
		m23: det * -(m.m00*A1312 - m.m01*A0312 + m.m03*A0112),
		m30: det * -(m.m10*A1223 - m.m11*A0223 + m.m12*A0123),
		m31: det * (m.m00*A1223 - m.m01*A0223 + m.m02*A0123),
		m32: det * -(m.m00*A1213 - m.m01*A0213 + m.m02*A0113),
		m33: det * (m.m00*A1212 - m.m01*A0212 + m.m02*A0112),
	}

}

func Up() {
	camPos.y += step
}

func Down() {
	camPos.y -= step
}

func StrafeRight() {
	right := normalize(cross(vec{0, 1, 0}, camDir))
	camPos = add(camPos, right)
}

func StrafeLeft() {
	right := normalize(cross(vec{0, 1, 0}, camDir))
	camPos = sub(camPos, right)
}

func TurnRight() {
	theta := math.Atan(camDir.x / camDir.z)
	theta += 0.01
	camDir.x = math.Sin(theta)
	camDir.z = math.Cos(theta)
}

func TurnLeft() {
	theta := math.Atan(camDir.x / camDir.z)
	theta -= 0.01
	camDir.x = math.Sin(theta)
	camDir.z = math.Cos(theta)
}

func Forward() {
	camPos = add(camPos, camDir)
}

func Back() {
	camPos = sub(camPos, camDir)
}

var step = 0.1
