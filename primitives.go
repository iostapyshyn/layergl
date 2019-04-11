package layergl

import (
	"math"
)

type Color struct {
	R, G, B, A float64
}

type Point struct {
	X, Y float64
}

type VertexObject struct {
	Vertices []Point
	Indices  []int
}

func Rectangle(c Point, d1, d2 float64) (polygon *VertexObject) {
	polygon = new(VertexObject)
	polygon.Vertices = []Point{
		{c.X, c.Y},
		{c.X, c.Y + d2},
		{c.X + d1, c.Y},
		{c.X + d1, c.Y + d2},
	}
	polygon.Indices = []int{
		0, 1, 2,
		1, 2, 3,
	}

	return polygon
}

func Square(c Point, d float64) (polygon *VertexObject) {
	polygon = new(VertexObject)
	polygon.Vertices = []Point{
		{c.X - d/2, c.Y - d/2},
		{c.X - d/2, c.Y + d/2},
		{c.X + d/2, c.Y - d/2},
		{c.X + d/2, c.Y + d/2},
	}
	polygon.Indices = []int{
		0, 1, 2,
		1, 2, 3,
	}

	return polygon
}

func Line(points ...Point) (polygon *VertexObject) {
	polygon = new(VertexObject)
	for i, v := range points {
		polygon.Vertices = append(polygon.Vertices, v)
		polygon.Indices = append(polygon.Indices, i)
	}

	return polygon
}

func Triangles(vertices ...Point) (polygon *VertexObject) {
	polygon = new(VertexObject)
	for i := 0; i+3 <= len(vertices); i += 3 {
		polygon.Vertices = append(polygon.Vertices, vertices[i], vertices[i+1], vertices[i+2])
		polygon.Indices = append(polygon.Indices, i, i+1, i+2)
	}

	return polygon
}

func (p *VertexObject) getVertexArray() (va []float32, elements []uint32) {
	for _, v := range p.Vertices {
		va = append(va, float32(v.X))
		va = append(va, float32(v.Y))
	}

	for _, i := range p.Indices {
		elements = append(elements, uint32(i))
	}

	return
}

func (p *VertexObject) Move(x, y float64) {
	for i, _ := range p.Vertices {
		p.Vertices[i].X += x
		p.Vertices[i].Y += y
	}
}

func (p *VertexObject) CenterAt(point Point) {
	center := p.Centroid()
	p.Move(-center.X+point.X, -center.Y+point.Y)
}

func (v *VertexObject) RotateDeg(angle float64) {
	v.RotateRad(angle * math.Pi / 180)
}

func (v *VertexObject) RotateRad(angle float64) {
	center := v.Centroid()
	for i, _ := range v.Vertices {
		pX := v.Vertices[i].X - center.X
		pY := v.Vertices[i].Y - center.Y
		v.Vertices[i].X = pX*math.Cos(angle) - pY*math.Sin(angle)
		v.Vertices[i].Y = pX*math.Sin(angle) + pY*math.Cos(angle)
		v.Vertices[i].X += center.X
		v.Vertices[i].Y += center.Y
	}
}

func (v *VertexObject) Scale(scale float64) {
	center := v.Centroid()
	for i, _ := range v.Vertices {
		v.Vertices[i].X -= center.X
		v.Vertices[i].Y -= center.Y
		v.Vertices[i].X *= scale
		v.Vertices[i].Y *= scale
		v.Vertices[i].X += center.X
		v.Vertices[i].Y += center.Y
	}
}
