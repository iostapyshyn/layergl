package layergl

import "math"

type Color = [4]float32

type Point struct {
	X, Y float32
}

type VertexObject struct {
	Vertices []Point
	Indices  []uint32
}

func Rectangle(c Point, d1, d2 float32) (polygon VertexObject) {
	polygon.Vertices = []Point{
		{c.X, c.Y},
		{c.X, c.Y + d2},
		{c.X + d1, c.Y},
		{c.X + d1, c.Y + d2},
	}
	polygon.Indices = []uint32{
		0, 1, 2,
		1, 2, 3,
	}
	return polygon
}

func Square(c Point, d float32) (polygon VertexObject) {
	polygon.Vertices = []Point{
		{c.X - d/2, c.Y - d/2},
		{c.X - d/2, c.Y + d/2},
		{c.X + d/2, c.Y - d/2},
		{c.X + d/2, c.Y + d/2},
	}
	polygon.Indices = []uint32{
		0, 1, 2,
		1, 2, 3,
	}
	return
}

func Line(points ...Point) (polygon VertexObject) {
	for i, v := range points {
		polygon.Vertices = append(polygon.Vertices, v)
		polygon.Indices = append(polygon.Indices, uint32(i))
	}
	return polygon
}

func Triangles(vertices ...Point) (polygon VertexObject) {
	for i := 0; i+3 <= len(vertices); i += 3 {
		polygon.Vertices = append(polygon.Vertices,
			vertices[i], vertices[i+1], vertices[i+2])
		polygon.Indices = append(polygon.Indices,
			uint32(i), uint32(i+1), uint32(i+2))
	}
	return polygon
}

func (p VertexObject) GetVertexArray() (va []float32, elements []uint32) {
	for _, v := range p.Vertices {
		va = append(va, v.X)
		va = append(va, v.Y)
	}
	return va, p.Indices
}

func (v VertexObject) XMostPoint() (point Point) {
	point = Point{float32(math.Inf(-1)), float32(math.Inf(-1))}
	for _, p := range v.Vertices {
		if p.X > point.X {
			point = p
		}
	}
	return point
}

func (v VertexObject) YMostPoint() (point Point) {
	point = Point{float32(math.Inf(-1)), float32(math.Inf(-1))}
	for _, p := range v.Vertices {
		if p.Y > point.Y {
			point = p
		}
	}
	return point
}

func (v VertexObject) XLeastPoint() (point Point) {
	point = Point{float32(math.Inf(+1)), float32(math.Inf(+1))}
	for _, p := range v.Vertices {
		if p.X < point.X {
			point = p
		}
	}
	return point
}

func (v VertexObject) YLeastPoint() (point Point) {
	point = Point{float32(math.Inf(+1)), float32(math.Inf(+1))}
	for _, p := range v.Vertices {
		if p.Y < point.Y {
			point = p
		}
	}
	return point
}

func (v VertexObject) Centroid() (center Point) {
	for _, p := range v.Vertices {
		center.X += p.X
		center.Y += p.Y
	}

	center.X /= float32(len(v.Vertices))
	center.Y /= float32(len(v.Vertices))
	return center
}

func (p VertexObject) Move(x, y float32) (new VertexObject) {
	new = p
	new.Vertices = make([]Point, len(p.Vertices))
	for i, v := range p.Vertices {
		new.Vertices[i].X = v.X + x
		new.Vertices[i].Y = v.Y + y
	}
	return new
}

func (p VertexObject) CenterAt(point Point) (new VertexObject) {
	center := p.Centroid()
	return p.Move(-center.X+point.X, -center.Y+point.Y)
}

func (v VertexObject) RotateDeg(angle float64) VertexObject {
	return v.RotateRad(angle * math.Pi / 180)
}

func (v VertexObject) RotateRad(angle float64) (new VertexObject) {
	new = v
	center := v.Centroid()
	new.Vertices = make([]Point, len(v.Vertices))
	for i, p := range v.Vertices {
		pX := p.X - center.X
		pY := p.Y - center.Y
		new.Vertices[i].X = pX*float32(math.Cos(angle)) - pY*float32(math.Sin(angle))
		new.Vertices[i].Y = pX*float32(math.Sin(angle)) + pY*float32(math.Cos(angle))
		new.Vertices[i].X += center.X
		new.Vertices[i].Y += center.Y
	}
	return new
}

func (v VertexObject) Scale(scale float64) (new VertexObject) {
	new = v
	center := new.Centroid()
	new.Vertices = make([]Point, len(v.Vertices))
	for i, p := range v.Vertices {
		new.Vertices[i].X = p.X - center.X
		new.Vertices[i].Y = p.Y - center.Y
		new.Vertices[i].X *= float32(scale)
		new.Vertices[i].Y *= float32(scale)
		new.Vertices[i].X += center.X
		new.Vertices[i].Y += center.Y
	}
	return new
}
