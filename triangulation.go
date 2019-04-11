package layergl

import (
	"fmt"
)

// Signed area of a triangle.
func sign(a, b, c Point) float64 {
	return (a.X-c.X)*(b.Y-c.Y) - (b.X-c.X)*(a.Y-c.Y)
}

// Is point P inside a triangle ABC?
func triangleContains(P, A, B, C Point) bool {

	if P == A || P == B || P == C {
		return false
	}

	area1 := sign(P, A, B)
	area2 := sign(P, B, C)
	area3 := sign(P, C, A)

	hasNegative := area1 < 0 || area2 < 0 || area3 < 0
	hasPositive := area1 > 0 || area2 > 0 || area3 > 0

	return !(hasNegative && hasPositive)

}

// Vertices have to be in clockwise order
func findEar(vert []Point, indexes []int) int {
	if len(indexes) < 3 {
		return 0
	}

	if len(indexes) == 3 {
		return 1
	}

	var success bool
	for i := 0; i < len(indexes)-2; i++ {
		if sign(vert[indexes[i]], vert[indexes[i+1]], vert[indexes[i+2]]) < 0 { // if P(i+1) vertex makes right turn...
			success = true
			for _, p := range vert { // Check if there are no other vertices inside of a triangle.
				if triangleContains(p, vert[indexes[i]], vert[indexes[i+1]], vert[indexes[i+2]]) {
					success = false
				}
			}
		}

		// If not, we found an ear.
		if success {
			return i + 1
		}
	}

	// Function returns 0 if there are no ears to be found or when given invalid input
	return 0
}

func (vo *VertexObject) Triangulate() error {
	if len(vo.Vertices) < 3 {
		return fmt.Errorf("unable perform triangulation of the polygon")
	}

	indexes := make([]int, len(vo.Vertices))
	for i := 0; i < len(vo.Vertices); i++ {
		indexes[i] = i
	}

	// In most cases triangulation of n vertices creates n-2 triangles.
	vo.Indices = make([]int, 0, (len(vo.Vertices)-2)*3)
	for len(indexes) >= 3 {
		if i := findEar(vo.Vertices, indexes); i == 0 {
			vo.Indices = vo.Indices[:0]
			return fmt.Errorf("unable perform triangulation of the polygon")
		} else {
			vo.Indices = append(vo.Indices, indexes[i-1], indexes[i], indexes[i+1])
			indexes = append(indexes[:i], indexes[i+1:]...)
		}
	}

	return nil
}

func PolygonFromVertices(p ...Point) (*VertexObject, error) {
	vo := new(VertexObject)
	vo.Vertices = p

	if err := vo.Triangulate(); err != nil {
		return vo, err
	}

	return vo, nil
}
