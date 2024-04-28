package utils

type BoundingBox struct {
  width float64
  height float64
}

type Collider struct {
  Pos *Vector2DFloat64
  box BoundingBox
}

func checkOverlap(x1min, x1max, x2min, x2max float64) bool {
  if x1max > x2min && x2max > x1min {
    return true
  }
  return false
}

func (c *Collider) CheckCollision(c2 *Collider) bool {
  x1min, x1max := c.Pos.X - c.box.width/2, c.Pos.X + c.box.width/2
  y1min, y1max := c.Pos.Y - c.box.height/2, c.Pos.Y + c.box.height/2
  x2min, x2max := c2.Pos.X - c2.box.width/2, c2.Pos.X + c2.box.width/2
  y2min, y2max := c2.Pos.Y - c2.box.height/2, c2.Pos.Y + c2.box.height/2
  if checkOverlap(x1min, x1max, x2min, x2max) && checkOverlap(y1min, y1max, y2min, y2max) {
    return true
  }
  return false
}

func NewCollider(pos *Vector2DFloat64, w, h float64) *Collider {
  return &Collider{
    Pos: pos,
    box: BoundingBox{width: w, height: h},
  }
}
