package geometry

import (
	"testing"

	"github.com/faiface/pixel"
	"github.com/stretchr/testify/assert"
)

func TestLineSegmentIntersect(t *testing.T) {
	assert.True(t, LineSegmentsIntersect(
		pixel.L(pixel.V(0, -1), pixel.V(0, 1)),
		pixel.L(pixel.V(-1, 0), pixel.V(1, 0)),
	))
	assert.True(t, LineSegmentsIntersect(
		pixel.L(pixel.V(0, -200), pixel.V(0, -50)),
		pixel.L(pixel.V(-92, -163), pixel.V(60, -162)),
	))
}
