package utils

type Vector2DInt64 struct {
  X int64
  Y int64
}

type Vector2DFloat64 struct {
  X float64
  Y float64
}

func (v1 *Vector2DInt64) Add(v2 *Vector2DInt64) {
  v1.X += v2.X
  v1.Y += v2.Y
}

func (v1 *Vector2DFloat64) Add(v2 *Vector2DFloat64) {
  v1.X += v2.X
  v1.Y += v2.Y
}

func (v1 *Vector2DInt64) Sub(v2 *Vector2DInt64) {
  v1.X -= v2.X
  v1.Y -= v2.Y
}

func (v1 *Vector2DFloat64) Sub(v2 *Vector2DFloat64) {
  v1.X -= v2.X
  v1.Y -= v2.Y
}
