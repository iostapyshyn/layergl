package layergl

import "math"

// Return distance between two points.
func Distance(a, b Point) float64 {
	return math.Sqrt((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y))
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
