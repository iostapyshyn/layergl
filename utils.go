package layergl

import "math"

func Distance(a, b Point) float64 {
	return math.Sqrt((b.X-a.X)*(b.X-a.X) + (b.Y-a.Y)*(b.Y-a.Y))
}

func (v *VertexObject) XMostPoint() (point Point) {
	point = Point{math.Inf(-1), math.Inf(-1)}

	for _, p := range v.Vertices {
		if p.X > point.X {
			point = p
		}
	}

	return
}

func (v *VertexObject) YMostPoint() (point Point) {
	point = Point{math.Inf(-1), math.Inf(-1)}

	for _, p := range v.Vertices {
		if p.Y > point.Y {
			point = p
		}
	}

	return
}

func (v *VertexObject) XLeastPoint() (point Point) {
	point = Point{math.Inf(+1), math.Inf(+1)}
	for _, p := range v.Vertices {
		if p.X < point.X {
			point = p
		}
	}
	return point
}

func (v *VertexObject) YLeastPoint() (point Point) {
	point = Point{math.Inf(+1), math.Inf(+1)}
	for _, p := range v.Vertices {
		if p.Y < point.Y {
			point = p
		}
	}
	return point
}

func (v *VertexObject) Centroid() (center Point) {
	for _, p := range v.Vertices {
		center.X += p.X
		center.Y += p.Y
	}

	center.X /= float64(len(v.Vertices))
	center.Y /= float64(len(v.Vertices))
	return center
}
