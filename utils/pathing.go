package utils

import (
	"math"
)

type PathingFunction func([][]uint8, int, int, int, int, int, int, int, int, bool, int) (int, Vector2DInt64)

func dist[TA int | int64, TB int | int64](x1 TA, y1 TA, x2 TB, y2 TB) float64 {
  return math.Sqrt(math.Pow(float64(x1) - float64(x2), 2) + math.Pow(float64(y1) - float64(y2), 2))
}

func homePathing(matrix [][]uint8, cur_x, cur_y, _, home_x, home_y int) (int, Vector2DInt64){
  dir_x := [4]int64{-1, 1, 0, 0}
  dir_y := [4]int64{0, 0, -1, 1}

  pq := NewPriorityQueue[Vector2DInt64]()
  mp := make(map[Vector2DInt64]Vector2DInt64)

  pq.PushHeap(&PriorityQueueItem[Vector2DInt64]{
    Value: Vector2DInt64{X: int64(cur_x), Y: int64(cur_y)},
    Index: 0,
    Priority: 0,
  })

  mp[Vector2DInt64{X: int64(cur_x), Y: int64(cur_y)}] = Vector2DInt64{X: int64(-1), Y: int64(-1)}

  var mat_exp [29][23]int

  mat_exp[cur_y][cur_x] = 5

  for pq.Len() > 0 {
    pos := pq.PopHeap()
    flag := false

    for i := range 4 {
      n_pos := Vector2DInt64{X: pos.Value.X+dir_x[i], Y: pos.Value.Y+dir_y[i]}

      if n_pos.Y < 0 || n_pos.Y >= int64(len(matrix)) || n_pos.X < 0 || n_pos.X >= int64(len(matrix[0])) {continue}
      if matrix[n_pos.Y][n_pos.X] == 1 {continue}
      if _, ok := mp[n_pos]; ok {continue}

      pq.PushHeap(&PriorityQueueItem[Vector2DInt64]{
        Value: n_pos,
        Index: 0,
        Priority: pos.Priority+1,
      })
      mp[n_pos] = pos.Value

      mat_exp[n_pos.Y][n_pos.X] = i+1

      if n_pos.X == int64(home_x) && n_pos.Y == int64(home_y) {flag=true;break;}
    }
    if flag {break}
  }


  // Print the representation matrix to check validity
  //for i := range(29) {
  //  for j := range(23) {
  //    if mat_exp[i][j] == 2 { fmt.Print("←"); continue; }
  //    if mat_exp[i][j] == 1 { fmt.Print("→"); continue; }
  //    if mat_exp[i][j] == 4 { fmt.Print("↑"); continue; }
  //    if mat_exp[i][j] == 3 { fmt.Print("↓"); continue; }
  //    if mat_exp[i][j] == 5 { fmt.Print("*"); continue; }
  //    fmt.Print(".")
  //  }
  //  fmt.Println()
  //}
  //fmt.Println()
  // End of printing the representation matrix

  prev := Vector2DInt64{X: int64(home_x), Y: int64(home_y)}
  next := mp[prev]

  path := make(map[Vector2DInt64]Vector2DInt64)
  for next.X != -1 && next.Y != -1 {path[next] = prev; prev, next = next, mp[next];}

  prev = Vector2DInt64{X: int64(cur_x), Y: int64(cur_y)}
  next = path[prev]

  n_dir := Vector2DInt64{X: next.X-prev.X, Y: next.Y-prev.Y}
  i := Vector2DInt64{X: next.X - prev.X, Y: next.Y - prev.Y}
  for i == n_dir {
    prev, next = next, path[next]
    i = Vector2DInt64{X: next.X - prev.X, Y: next.Y - prev.Y}
  }

  dir := 0
  switch n_dir {
  case Vector2DInt64{X: -1, Y: 0}:
    dir = 0
  case Vector2DInt64{X: 1,  Y: 0}:
    dir = 1
  case Vector2DInt64{X: 0, Y: -1}:
    dir = 2
  case Vector2DInt64{X: 0, Y: 1}:
    dir = 3
  }
  return dir, prev
}

func pathing(matrix [][]uint8, cur_x, cur_y, ghost_dir, tar_x, tar_y int) (int, Vector2DInt64) {
  anti_ghost_dir := 0
  if ghost_dir == 0 {anti_ghost_dir=1}
  if ghost_dir == 1 {anti_ghost_dir=0}
  if ghost_dir == 2 {anti_ghost_dir=3}
  if ghost_dir == 3 {anti_ghost_dir=2}
  
  dir_x := [4]int64{-1, 1, 0, 0}
  dir_y := [4]int64{0, 0, -1, 1}

  cur_pos := Vector2DInt64{X: int64(cur_x), Y: int64(cur_y)}
  n_pos := Vector2DInt64{X: 10000000, Y: 10000000}
  n_dir := 0

  for i := range(4) {
    if i == anti_ghost_dir {continue}
    pos := Vector2DInt64{X: cur_pos.X+dir_x[i], Y: cur_pos.Y+dir_y[i]}

    if pos.Y < 0 {pos.Y = int64(len(matrix)-1)}
    if pos.Y >= int64(len(matrix)) {pos.Y = 0}
    if pos.X < 0 {pos.X = int64(len(matrix[0])-1)}
    if pos.X >= int64(len(matrix[0])) {pos.X = 0}

    if matrix[pos.Y][pos.X] == 1 {continue}

    if dist(pos.X, pos.Y, tar_x, tar_y) < dist(n_pos.X, n_pos.Y, tar_x, tar_y) {
      n_pos = pos
      n_dir = i
    }
  }

  if n_dir == 0 {anti_ghost_dir=1}
  if n_dir == 1 {anti_ghost_dir=0}
  if n_dir == 2 {anti_ghost_dir=3}
  if n_dir == 3 {anti_ghost_dir=2}

  if n_pos.X == 10000000 || n_pos.Y == 10000000 {
    n_pos = Vector2DInt64{X: cur_pos.X+dir_x[anti_ghost_dir], Y: cur_pos.Y+dir_y[anti_ghost_dir]}
    n_dir, anti_ghost_dir = anti_ghost_dir, n_dir
  }

  for cnt_dirs := (1<<n_dir); cnt_dirs == (1<<n_dir); {
    cnt_dirs = 0
    dir_pos := n_pos
    for i := range(4) {
      if i == anti_ghost_dir {continue}

      pos := Vector2DInt64{X: n_pos.X+dir_x[i], Y: n_pos.Y+dir_y[i]}

      if pos.Y < 0 {pos.Y = int64(len(matrix)-1)}
      if pos.Y >= int64(len(matrix)) {pos.Y = 0}
      if pos.X < 0 {pos.X = int64(len(matrix[0])-1)}
      if pos.X >= int64(len(matrix[0])) {pos.X = 0}

      if matrix[pos.Y][pos.X] == 1 {continue}

      cnt_dirs = (cnt_dirs | (1<<i))
      dir_pos = pos
    }
    if cnt_dirs == (1<<n_dir) { n_pos = dir_pos }
  }

  return n_dir, n_pos
}

func RedGhostPathing(matrix [][]uint8, cur_x, cur_y, ghost_dir, pac_x, pac_y, _, _, _ int, mode bool, state int) (int, Vector2DInt64) {
  if state == 0 {return homePathing(matrix, cur_x, cur_y, ghost_dir, pac_x, pac_y)}
  if mode && state == 2 {return pathing(matrix, cur_x, cur_y, ghost_dir, pac_x, pac_y)}
  return pathing(matrix, cur_x, cur_y, ghost_dir, len(matrix[0]), -10)
}

func PinkGhostPathing(matrix [][]uint8, cur_x, cur_y, ghost_dir, pac_x, pac_y, pac_dir, _, _ int, mode bool, state int) (int, Vector2DInt64) {
  dir_x := [4]int{0, 1, 0, -1}
  dir_y := [4]int{-1, 0, 1, 0}
  if state == 0 {return homePathing(matrix, cur_x, cur_y, ghost_dir, pac_x, pac_y)}
  if mode && state == 2 {return pathing(matrix, cur_x, cur_y, ghost_dir, pac_x+4*dir_x[pac_dir], pac_y+4*dir_y[pac_dir])}
  return pathing(matrix, cur_x, cur_y, ghost_dir, 0, -10)
}

func BlueGhostPathing(matrix [][]uint8, cur_x, cur_y, ghost_dir, pac_x, pac_y, pac_dir, red_ghost_x, red_ghost_y int, mode bool, state int) (int, Vector2DInt64) {
  dir_x := [4]int{0, 1, 0, -1}
  dir_y := [4]int{-1, 0, 1, 0}
  n_pac_x, n_pac_y := pac_x+2*dir_x[pac_dir], pac_y+2*dir_y[pac_dir]

  delta_x, delta_y := n_pac_x - red_ghost_x, n_pac_y-red_ghost_y
  if state == 0 {return homePathing(matrix, cur_x, cur_y, ghost_dir, pac_x, pac_y)}
  if mode && state == 2 {return pathing(matrix, cur_x, cur_y, ghost_dir, n_pac_x+delta_x, n_pac_y+delta_y)}
  return pathing(matrix, cur_x, cur_y, ghost_dir, len(matrix[0]), len(matrix)+10)
}

func OrangeGhostPathing(matrix [][]uint8, cur_x, cur_y, ghost_dir, pac_x, pac_y, _, _, _ int, mode bool, state int) (int, Vector2DInt64) {
  if state == 0 {return homePathing(matrix, cur_x, cur_y, ghost_dir, pac_x, pac_y)}
  if mode && state == 2{
    if dist(cur_x, cur_y, pac_x, pac_y) > 4 {return pathing(matrix, cur_x, cur_y, ghost_dir, pac_x, pac_y)}
  }
  return pathing(matrix, cur_x, cur_y, ghost_dir, 0, len(matrix)+10)
}
