package utils

func Key2D[T int64 | int](i, j T) Vector2DInt64 {
  return Vector2DInt64{X: int64(j), Y: int64(i)}
}

func Key2DToXY(pos Vector2DInt64) (int, int) {
  return int(pos.X), int(pos.Y)
}
