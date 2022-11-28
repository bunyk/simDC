package ui

// FASTER LINE SEGMENT INTERSECTION
// GRAPHICS GEMS III
func lineSegmentsIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	ax := x2 - x1
	ay := y2 - y1
	bx := x3 - x4
	by := y3 - y4

	d := ay*bx - ax*by
	if d == 0 { // collinear
		return false
	}
	cx := x1 - x3
	cy := y1 - y3
	a := by*cx - bx*cy
	if d > 0 {
		if a < 0 || a > d {
			return false
		}
	} else {
		if a > 0 || a < d {
			return false
		}
	}
	b := ax*cy - ay*cx
	if d > 0 {
		if b < 0 || b > d {
			return false
		}
	} else {
		if b > 0 || b < d {
			return false
		}
	}
	return true
}

func sqr(x float64) float64 {
	return x * x
}

// https://math.stackexchange.com/a/4088608/571313
func lineCircleIntersect(x1, y1, x2, y2, cx, cy, r float64) bool {
	// Change coordinate origin to be at the circle center
	x1 = x1 - cx
	y1 = y1 - cy
	x2 = x2 - cx
	y2 = y2 - cy

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

func lineRectangleIntersect(x1, y1, x2, y2, xmin, ymin, xmax, ymax float64) bool {
	return (lineSegmentsIntersect(x1, y1, x2, y2, xmin, ymin, xmax, ymin) || // bottom
		lineSegmentsIntersect(x1, y1, x2, y2, xmax, ymin, xmax, ymax) || // right
		lineSegmentsIntersect(x1, y1, x2, y2, xmin, ymin, xmin, ymax) || // left
		lineSegmentsIntersect(x1, y1, x2, y2, xmin, ymax, xmax, ymax)) // top
}
