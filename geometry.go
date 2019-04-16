package layergl

import (
	"log"
	"math"
)

type Color struct {
	R, G, B, A float64
}

type Point struct {
	X, Y float64
}

type Rect struct {
	X1, Y1 float64
	X2, Y2 float64
}

type VertexObject struct {
	Vertices []Point
	Indices  []int
}

// Return distance between two points.
func Distance(a, b Point) float64 {
	return math.Sqrt((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y))
}

// Returns geometrical center, center of mass (centroid) of a polygon.
func (v VertexObject) Centroid() (center Point) {
	for _, p := range v.Vertices {
		center.X += p.X
		center.Y += p.Y
	}

	center.X /= float64(len(v.Vertices))
	center.Y /= float64(len(v.Vertices))
	return center
}

func Rectangle(rect Rect) (polygon *VertexObject) {
	polygon = new(VertexObject)
	polygon.Vertices = []Point{
		{rect.X1, rect.Y1},
		{rect.X1, rect.Y2},
		{rect.X2, rect.Y1},
		{rect.X2, rect.Y2},
	}
	polygon.Indices = []int{
		0, 1, 2,
		1, 2, 3,
	}

	return polygon
}

func Triangles(vertices []Point) (polygon *VertexObject) {
	if len(vertices)%3 != 0 {
		log.Println("Number of vertices for layergl.Triangle() should be multiple of 3.")
	}

	polygon = new(VertexObject)
	for i := 0; i+3 <= len(vertices); i += 3 {
		polygon.Vertices = append(polygon.Vertices, vertices[i], vertices[i+1], vertices[i+2])
		polygon.Indices = append(polygon.Indices, i, i+1, i+2)
	}

	return polygon
}

// Returns vertex array of Rect in a proper format to load into the buffers.
func (r Rect) vertexArray() ([]float32, []uint32) {
	va := []float32{
		float32(r.X1), float32(r.Y1),
		float32(r.X1), float32(r.Y2),
		float32(r.X2), float32(r.Y1),
		float32(r.X2), float32(r.Y2),
	}

	elements := []uint32{
		0, 1, 2,
		1, 2, 3,
	}

	return va, elements
}

// Returns vertex array of VertexObject in a proper format to load into the buffers.
func (p VertexObject) vertexArray() (va []float32, elements []uint32) {
	for _, v := range p.Vertices {
		va = append(va, float32(v.X))
		va = append(va, float32(v.Y))
	}

	for _, i := range p.Indices {
		elements = append(elements, uint32(i))
	}

	return
}

// Translation of the VertexObject by x, y pixels in the respective directions.
func (p *VertexObject) Move(x, y float64) {
	for i, _ := range p.Vertices {
		p.Vertices[i].X += x
		p.Vertices[i].Y += y
	}
}

// Centers VertexObject at exact point.
func (p *VertexObject) CenterAt(point Point) {
	center := p.Centroid()
	p.Move(-center.X+point.X, -center.Y+point.Y)
}

// Rotation.
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

// Scaling.
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

// Return rectangle approximating bounds of the polygon.
func (v VertexObject) Bounds() Rect {
	if !(len(v.Vertices) > 0) {
		return Rect{}
	}

	xmost := math.Inf(-1)
	ymost := math.Inf(-1)
	xleast := math.Inf(+1)
	yleast := math.Inf(+1)

	for _, p := range v.Vertices {
		if p.X > xmost {
			xmost = p.X
		}
		if p.Y > ymost {
			ymost = p.Y
		}
		if p.X < xleast {
			xleast = p.X
		}
		if p.Y < yleast {
			yleast = p.Y
		}
	}

	return Rect{xleast, yleast, xmost, ymost}
}

func (rect Rect) Contains(p Point) bool {
	return p.X > rect.X1 && p.X < rect.X2 && p.Y > rect.Y1 && p.Y < rect.Y2
}

func (rect Rect) Intersects(r2 Rect) bool {
	return rect.Contains(Point{r2.X1, r2.Y1}) || rect.Contains(Point{r2.X2, r2.Y2}) ||
		rect.Contains(Point{r2.X1, r2.Y2}) || rect.Contains(Point{r2.X2, r2.Y1})
}

func (rect Rect) Height() float64 {
	return rect.Y2 - rect.Y1
}

func (rect Rect) Width() float64 {
	return rect.X2 - rect.X1
}

func (rect Rect) Area() float64 {
	return rect.Height() * rect.Width()
}
