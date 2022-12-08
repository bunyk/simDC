package geometry

import (
	"fmt"

	"github.com/faiface/pixel"
)

// FASTER LINE SEGMENT INTERSECTION
// GRAPHICS GEMS III
func LineSegmentsIntersect(a, b pixel.Line) bool {
	da := a.B.Sub(a.A)
	db := b.A.Sub(b.B)

	d := da.Y*db.X - da.X*db.Y
	if d == 0 { // collinear
		return false
	}
	c := a.A.Sub(b.A)
	d1 := db.Y*c.X - db.X*c.Y
	fmt.Println(c, d1)
	if d > 0 {
		if d1 < 0 || d1 > d {
			return false
		}
	} else {
		if d1 > 0 || d1 < d {
			return false
		}
	}
	d2 := da.X*c.Y - da.Y*c.X
	if d > 0 {
		if d2 < 0 || d2 > d {
			return false
		}
	} else {
		if d2 > 0 || d2 < d {
			return false
		}
	}
	return true
}

func sqr(x float64) float64 {
	return x * x
}

// https://math.stackexchange.com/a/4088608/571313
func LineCircleIntersect(l pixel.Line, center pixel.Vec, r float64) bool {
	// Change coordinate origin to be at the circle center
	x1 := l.A.X - center.X
	y1 := l.A.Y - center.Y
	x2 := l.B.X - center.X
	y2 := l.B.Y - center.Y

	// Square the radius to avoid needing any square roots
	r2 := r * r

	// Check if endpoints are in the circle
	if sqr(x1)+sqr(y1) <= r2 {
		return true
	}
	if sqr(x2)+sqr(y2) <= r2 {
		return true
	}
	// Calculate the line segment length squared
	len2 := sqr(x1-x2) + sqr(y1-y2)
	if len2 == 0 {
		return false
	}

	// find the perpendicular vector to the line segment
	nx := y2 - y1
	ny := x1 - x2

	// find the distance squared from center to line, times len2
	dist2 := sqr(nx*x1 + ny*y1)

	// Check if full line intersects the circle. If not, return false
	if dist2 > len2*r2 {
		return false
	}

	// Calculate the distance from (x1,y1) to the point of closest approach) times the segment length
	index := (x1*(x1-x2) + y1*(y1-y2))

	// Check if point of closest approach is inside the segment.
	if index < 0 {
		return false
	}
	if index > len2 {
		return false
	}
	return true
}

func LineRectangleIntersect(l pixel.Line, r pixel.Rect) bool {
	for _, e := range r.Edges() {
		if LineSegmentsIntersect(l, e) {
			return true
		}
	}
	return false
}
