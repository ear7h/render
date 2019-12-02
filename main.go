package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

type Color [4]float64

func (c Color) Add(c1 Color) Color {
	return Color{
		c[0] + c1[0],
		c[1] + c1[1],
		c[2] + c1[2],
		c[3] + c1[3],
	}
}

func (c Color) Mul(n float64) Color {
	return Color{
		c[0] * n,
		c[1] * n,
		c[2] * n,
		c[3] * n,
	}
}

/*
func (c Color) RGB [3]uint8 {
	return
}
*/

type Point [3]float64

func (p0 Point) Sub(p1 Point) Point {
	return Point{
		p0[0] - p1[0],
		p0[1] - p1[1],
		p0[2] - p1[2],
	}
}

func (p0 Point) Add(p1 Point) Point {
	return Point{
		p0[0] + p1[0],
		p0[1] + p1[1],
		p0[2] + p1[2],
	}
}

func (p Point) Norm() float64 {
	return math.Sqrt(p[0]*p[0] + p[1]*p[1] + p[2]*p[2])
}

func (p Point) Dot(p1 Point) float64 {
	return p[0]*p1[0] + p[1]*p1[1] + p[2]*p1[2]
}

func (p Point) Cross(v Point) Point {
	a := p
	b := v
	return Point{
		a[1]*b[2] - a[2]*b[1],
		a[2]*b[0] - a[0]*b[2],
		a[0]*b[1] - a[1]*b[0],
	}
}

func (p Point) Div(n float64) Point {
	return Point{
		p[0] / n,
		p[1] / n,
		p[2] / n,
	}
}

func (p Point) Mul(n float64) Point {
	return Point{
		p[0] * n,
		p[1] * n,
		p[2] * n,
	}
}

func (p Point) Unit() Point {
	return p.Div(p.Norm())
}

func randPoint() Point {
	return Point{
		rand.Float64(),
		rand.Float64(),
		rand.Float64(),
	}
}

func genLines() [][2]Point {
	points := make([]Point, 20)
	for i := range points {
		points[i] = rndSphere1(0.4)
	}

	lines := make([][2]Point, 200)
	for i := range lines {
		i1 := rand.Intn(len(points))
		i2 := rand.Intn(len(points))
		lines[i][0] = points[i1]
		lines[i][1] = points[i2]
	}

	return lines
}

type Camera struct {
	// The position of the center of the frame
	Position Point
	// Forward and top left
	Orientation   [2]Point
	Height, Width int
	Pixels        []Color
}

func (c Camera) StorePoint(pt Point, color Color) {
	v1 := pt.Sub(c.Position)
	if v1.Dot(c.Orientation[0]) < 0 {
		return
	}

	tl := c.Orientation[1]
	tr := c.Orientation[0].Cross(c.Orientation[1])

	yaxis := tl.Add(tr).Unit()
	xaxis := tl.Mul(-1).Add(tr).Unit()

	h := c.Height
	w := c.Width

	x := int((0.5 + xaxis.Dot(v1)) * float64(w))
	y := int((0.5 + yaxis.Dot(v1)) * float64(h))

	idx := x + y*w
	if idx < 0 || idx >= len(c.Pixels) {
		return
	}

	color0 := c.Pixels[idx]
	c.Pixels[idx] = color0.Add(color)
}

func (c Camera) Image() image.Image {
	img := image.NewRGBA(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Point{c.Width, c.Height},
		})

	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			pix := c.Pixels[y*c.Width+x].Mul(0xff)
			img.Set(x, y, color.NRGBA{
				R: uint8(math.Min(pix[0], 0xff)),
				G: uint8(math.Min(pix[1], 0xff)),
				B: uint8(math.Min(pix[2], 0xff)),
				A: uint8(math.Min(pix[3], 0xff)),
			})
		}
	}

	return img
}

func RenderLine(c Camera, l [2]Point, col Color) {
	const delta = 0.0001
	v := l[1].Sub(l[0])
	vn := v.Norm()
	v = v.Unit()

	c.StorePoint(l[0], col)

	for i := float64(0); i < vn; i += delta {
		c.StorePoint(l[0].Add(v.Mul(i)), col)
	}

	c.StorePoint(l[1], col)
}

func rndSphere(r float64) Point {
	rho := math.Pi * rand.Float64()
	theta := math.Pi * 2 * rand.Float64()
	r += 0.01 * rand.NormFloat64()

	return Point{
		r * math.Sin(theta) * math.Cos(rho),
		r * math.Sin(theta) * math.Sin(rho),
		r * math.Cos(theta),
	}
}

func rndSphere1(r float64) Point {
	rho := math.Pi * rand.Float64()
	theta := math.Pi * 2 * rand.Float64()

	return Point{
		r * math.Sin(theta) * math.Cos(rho),
		r * math.Sin(theta) * math.Sin(rho),
		r * math.Cos(theta),
	}
}

func RenderLineMC(c Camera, l [2]Point, col Color) {
	const delta = 0.0001
	v := l[1].Sub(l[0])
	vn := v.Norm()
	v = v.Unit()

	c.StorePoint(l[0], col)

	for i := float64(0); i < vn; i += delta {
		pt := l[0].Add(v.Mul(i))
		d := pt.Sub(c.Position).Norm()
		r := 0.2 * math.Pow(math.Abs(0.5-d), 2)
		pt = pt.Add(rndSphere(r))
		c.StorePoint(pt, col)
	}

	c.StorePoint(l[1], col)
}

func main() {
	dir := fmt.Sprintf("out/%d", time.Now().Unix())
	err := os.Mkdir(dir, 0777)
	if err != nil {
		log.Fatal(err)
	}

	lines := genLines()
	rot := NewRotationY(math.Pi/10)
	fmt.Println(rot)

	for i := 0; i < 20; i++ {
		camera := Camera{
			Position: Point{-0.5, 0, 0},
			Orientation: [2]Point{
				{1, 0, 0},
				{0, 0.5, 0.5},
			},
			Height: 300,
			Width:  300,
			Pixels: make([]Color, 300*300),
		}


		for i := range lines {
			if i == 0 {
				fmt.Println(lines[i][0])
			}

			lines[i][0] = rot.Apply(lines[i][0])
			lines[i][1] = rot.Apply(lines[i][1])

			if i == 0 {
				fmt.Println(lines[i][0])
			}

			RenderLineMC(camera, lines[i], Color{0.75, 0.75, 0.75, 0.01})
		}

		fname := fmt.Sprintf("%s/out%04d.png", dir, i)
		fmt.Println(fname)
		err := saveImage(fname, camera.Image())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func saveImage(name string, img image.Image) error {
	fd, err := os.OpenFile(
		name,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	return png.Encode(fd, img)
}
