package layergl

import (
	"fmt"
)

// Signed area of a triangle.
func sign(a, b, c Point) float64 {
	return (a.X-c.X)*(b.Y-c.Y) - (b.X-c.X)*(a.Y-c.Y)
}

// Returns true if point P is inside triangle ABC.
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

// Returns index of an ear to be cut from polygon.
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

	// findEar returns 0 if there are no ears to be found or when given invalid input.
	return 0
}

// Performs triangulation of VertexObject, writing to the Indices field.
func (vo *VertexObject) Triangulate() error {
	if len(vo.Vertices) < 3 {
		return fmt.Errorf("unable perform triangulation of the polygon")
	}

	// Stores list of ears still present in a polygon in a process of ear-cutting.
	ears := make([]int, len(vo.Vertices))
	for i := 0; i < len(vo.Vertices); i++ {
		ears[i] = i
	}

	// In most cases triangulation of n vertices creates n-2 triangles.
	vo.Indices = make([]int, 0, (len(vo.Vertices)-2)*3)
	for len(ears) >= 3 {
		if i := findEar(vo.Vertices, ears); i == 0 {
			// In case findEar fails.
			vo.Indices = vo.Indices[:0]
			return fmt.Errorf("unable perform triangulation of the polygon")
		} else {
			// Add new triangle.
			vo.Indices = append(vo.Indices, ears[i-1], ears[i], ears[i+1])

			// Cut the ear.
			ears = append(ears[:i], ears[i+1:]...)
		}
	}

	return nil
}

// Performs triangulation creating new VertexObject.
func FromVertices(p []Point) (*VertexObject, error) {
	vo := new(VertexObject)
	vo.Vertices = p

	if err := vo.Triangulate(); err != nil {
		return vo, err
	}

	return vo, nil
}
