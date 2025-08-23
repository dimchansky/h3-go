package c2go

// _v2dIntersect finds the intersection between two lines, assuming they
// intersect and not at endpoints. Ported from vec2d.c::_v2dIntersect
// Ported from H3 C: vec2d.c::v2dIntersect
func _v2dIntersect(p0, p1, p2, p3 *Vec2d) Vec2d {
	s1x := p1.X - p0.X
	s1y := p1.Y - p0.Y
	s2x := p3.X - p2.X
	s2y := p3.Y - p2.Y
	t := (s2x*(p0.Y-p2.Y) - s2y*(p0.X-p2.X)) / (-s2x*s1y + s1x*s2y)
	return Vec2d{X: p0.X + t*s1x, Y: p0.Y + t*s1y}
}
