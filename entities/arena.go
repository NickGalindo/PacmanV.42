package entities

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"time"

	"github.com/NickGalindo/Pacman/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
  arena_width = 23
  arena_height = 29
  arena_tile_size = 20

  ARENA_START_STATE = 0
  ARENA_PLAYING_STATE = 1
  ARENA_END_STATE = 2
  ARENA_RESET_STATE = 3
  ARENA_WIN_STATE = 4

  ARENA_SCORE_FONT_SIZE = 16
)

type Arena struct {
  Matrix [][]uint8
  Pos *utils.Vector2DFloat64
  Score int64
  w map[utils.Vector2DInt64]*Wall
  bc map[utils.Vector2DInt64]*BigCoin
  sc map[utils.Vector2DInt64]*SmallCoin
  pac *Pacman
  ghost_red *Ghost
  ghost_blue *Ghost
  ghost_pink *Ghost
  ghost_orange *Ghost
  start_button *StartButton
  start_text_box *TextBox
  end_text_box *TextBox
  win_text_box *TextBox
  state uint8
  arcadeFaceSource *text.GoTextFaceSource 
  totalPlayTime time.Duration
  lastPlayStart time.Time
}

func NewArena(pos *utils.Vector2DFloat64) (*Arena, error){
  pac_x, pac_y := (11.0*arena_tile_size)+pos.X, (17.0*arena_tile_size)+pos.Y
  pac, _ := NewPacman(&utils.Vector2DFloat64{X: pac_x, Y: pac_y})

  
  gho_x, gho_y := (9.0*arena_tile_size)+pos.X, (13.0*arena_tile_size)+pos.Y
  ghost_red, _ := NewGhost(&utils.Vector2DFloat64{X: gho_x, Y: gho_y}, "ghost_red.png", utils.RedGhostPathing)
  gho_x, gho_y = (13.0*arena_tile_size)+pos.X, (13.0*arena_tile_size)+pos.Y
  ghost_blue, _ := NewGhost(&utils.Vector2DFloat64{X: gho_x, Y: gho_y}, "ghost_blue.png", utils.BlueGhostPathing)
  gho_x, gho_y = (9.0*arena_tile_size)+pos.X, (15.0*arena_tile_size)+pos.Y
  ghost_pink, _ := NewGhost(&utils.Vector2DFloat64{X: gho_x, Y: gho_y}, "ghost_pink.png", utils.PinkGhostPathing)
  gho_x, gho_y = (13.0*arena_tile_size)+pos.X, (15.0*arena_tile_size)+pos.Y
  ghost_orange, _ := NewGhost(&utils.Vector2DFloat64{X: gho_x, Y: gho_y}, "ghost_orange.png", utils.OrangeGhostPathing)

  start_button, _ := NewStartButton(&utils.Vector2DFloat64{X: (pos.X-10+arena_tile_size*arena_width)/2, Y: (pos.Y+10+arena_tile_size*arena_height)/2})

  start_text_box, _ := NewTextBox(&utils.Vector2DFloat64{X: (pos.X-10+arena_tile_size*arena_width)/2, Y: (pos.Y+10+arena_tile_size*arena_height)/2}, "start_text_box.png")
  end_text_box, _ := NewTextBox(&utils.Vector2DFloat64{X: (pos.X-10+arena_tile_size*arena_width)/2, Y: (pos.Y+10+arena_tile_size*arena_height)/2}, "end_text_box.png")
  win_text_box, _ := NewTextBox(&utils.Vector2DFloat64{X: (pos.X-10+arena_tile_size*arena_width)/2, Y: (pos.Y+10+arena_tile_size*arena_height)/2}, "win_text_box.png")

  return &Arena{
    Matrix: [][]uint8{
      {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
      {1, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2, 1},
      {1, 3, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 3, 1},
      {1, 1, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1},
      {0, 0, 0, 0, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 0, 0, 0, 0},
      {0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0},
      {0, 0, 0, 0, 1, 3, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 3, 1, 0, 0, 0, 0},
      {1, 1, 1, 1, 1, 3, 1, 0, 1, 5, 0, 0, 0, 8, 1, 0, 1, 3, 1, 1, 1, 1, 1},
      {3, 3, 3, 3, 3, 3, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1, 3, 3, 3, 3, 3, 3},
      {1, 1, 1, 1, 1, 3, 1, 0, 1, 6, 0, 0, 0, 7, 1, 0, 1, 3, 1, 1, 1, 1, 1},
      {0, 0, 0, 0, 1, 3, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 3, 1, 0, 0, 0, 0},
      {0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 3, 1, 0, 0, 0, 0},
      {0, 0, 0, 0, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 0, 0, 0, 0},
      {1, 1, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 1, 1},
      {1, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 1, 3, 3, 3, 3, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 3, 1},
      {1, 3, 1, 1, 1, 3, 1, 1, 1, 1, 3, 1, 3, 1, 1, 1, 1, 3, 1, 1, 1, 3, 1},
      {1, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 1, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2, 1},
      {1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},

      //{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
      //{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
    },
    w: make(map[utils.Vector2DInt64]*Wall),
    bc: make(map[utils.Vector2DInt64]*BigCoin),
    sc: make(map[utils.Vector2DInt64]*SmallCoin),
    Pos: pos,
    pac: pac,
    ghost_red: ghost_red,
    ghost_blue: ghost_blue,
    ghost_pink: ghost_pink,
    ghost_orange: ghost_orange,
    start_button: start_button,
    end_text_box: end_text_box,
    win_text_box: win_text_box,
    start_text_box: start_text_box,
    state: ARENA_START_STATE,
    totalPlayTime: 0,
  }, nil
}

func (a *Arena) FindCellFromXY(x, y float64) (int, int) {
  return int(math.Round((x-a.Pos.X)/arena_tile_size)), int(math.Round((y-a.Pos.Y)/arena_tile_size))
}

func (a *Arena) FindXYFromCell(x, y int) (float64, float64) {
  return (float64(x)*arena_tile_size)+a.Pos.X, (float64(y)*arena_tile_size)+a.Pos.Y
}

func (a *Arena) HandlePressedKey(key ebiten.Key) {
  if a.state == ARENA_START_STATE {
    if key == ebiten.KeyQ {a.state = ARENA_RESET_STATE}
  }
  if a.state == ARENA_RESET_STATE {
    if key == ebiten.KeyS {a.lastPlayStart = time.Now(); a.state = ARENA_PLAYING_STATE}
  }
  if a.state == ARENA_PLAYING_STATE {
    a.pac.HandlePressedKey(key)
  }
  if a.state == ARENA_END_STATE || a.state == ARENA_WIN_STATE {
    if key == ebiten.KeyQ {os.Exit(0)}
  }
}

func (a *Arena) Init() {
  s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
  if err != nil {
    log.Fatal(err)
  }

  a.arcadeFaceSource = s

  a.ghost_red.CalculatePathPoint(a)
  a.ghost_blue.CalculatePathPoint(a)
  a.ghost_orange.CalculatePathPoint(a)
  a.ghost_pink.CalculatePathPoint(a)
}

func (a *Arena) Reset() {
  a.state = ARENA_RESET_STATE

  a.pac.state = PACMAN_ALIVE_STATE
  a.pac.Dir = PACMAN_DIR_RIGHT
  a.pac.last_dir = PACMAN_DIR_RIGHT

  a.pac.Pos.X, a.pac.Pos.Y = (11.0*arena_tile_size)+a.Pos.X, (17.0*arena_tile_size)+a.Pos.Y

  a.ghost_red.state = GHOST_ALIVE_STATE
  a.ghost_blue.state = GHOST_ALIVE_STATE
  a.ghost_pink.state = GHOST_ALIVE_STATE
  a.ghost_orange.state = GHOST_ALIVE_STATE

  a.ghost_red.mode = GHOST_ACTIVE_MODE
  a.ghost_blue.mode = GHOST_ACTIVE_MODE
  a.ghost_pink.mode = GHOST_ACTIVE_MODE
  a.ghost_orange.mode = GHOST_ACTIVE_MODE

  a.ghost_red.dir = GHOST_DIR_UP
  a.ghost_blue.dir = GHOST_DIR_UP
  a.ghost_pink.dir = GHOST_DIR_UP
  a.ghost_orange.dir = GHOST_DIR_UP

  a.ghost_red.mode_clock = time.Now()
  a.ghost_blue.mode_clock = time.Now()
  a.ghost_pink.mode_clock = time.Now()
  a.ghost_orange.mode_clock = time.Now()

  a.ghost_red.Pos.X, a.ghost_red.Pos.Y = a.ghost_red.home.X, a.ghost_red.home.Y
  a.ghost_blue.Pos.X, a.ghost_blue.Pos.Y = a.ghost_blue.home.X, a.ghost_blue.home.Y
  a.ghost_pink.Pos.X, a.ghost_pink.Pos.Y = a.ghost_pink.home.X, a.ghost_pink.home.Y
  a.ghost_orange.Pos.X, a.ghost_orange.Pos.Y = a.ghost_orange.home.X, a.ghost_orange.home.Y

  a.ghost_red.CalculatePathPoint(a)
  a.ghost_blue.CalculatePathPoint(a)
  a.ghost_orange.CalculatePathPoint(a)
  a.ghost_pink.CalculatePathPoint(a)
}

func (a *Arena) EndGame() {
  a.state = ARENA_END_STATE
}

func (a *Arena) GameUpdate() {
  if a.pac.state == PACMAN_DEAD_STATE {
    if a.pac.frame == -1 {
      a.totalPlayTime += time.Since(a.lastPlayStart)
      a.Reset()
      if a.pac.Lives <= 0 {
        a.EndGame()
      }
    }
    return
  }

  // Update pacman
  red_ghost_death, blue_ghost_death, pink_ghost_death, orange_ghost_death   := a.ghost_red.state, a.ghost_blue.state, a.ghost_pink.state, a.ghost_orange.state
  a.pac.Update(a.w, []*Ghost{a.ghost_red, a.ghost_blue, a.ghost_pink, a.ghost_orange}, &a.Score)

  if red_ghost_death != GHOST_DEATH_STATE && red_ghost_death != a.ghost_red.state {a.ghost_red.CalculatePathPoint(a)}
  if blue_ghost_death != GHOST_DEATH_STATE && blue_ghost_death != a.ghost_blue.state {a.ghost_blue.CalculatePathPoint(a)}
  if pink_ghost_death != GHOST_DEATH_STATE && pink_ghost_death != a.ghost_pink.state {a.ghost_pink.CalculatePathPoint(a)}
  if orange_ghost_death != GHOST_DEATH_STATE && orange_ghost_death != a.ghost_orange.state {a.ghost_orange.CalculatePathPoint(a)}

  if a.pac.Pos.X < a.Pos.X { a.pac.Pos.X = a.Pos.X+arena_width*arena_tile_size-1 }
  if a.pac.Pos.X >= a.Pos.X+arena_width*arena_tile_size { a.pac.Pos.X = a.Pos.X+1 }
  if a.pac.Pos.Y < a.Pos.Y { a.pac.Pos.Y = a.Pos.Y+arena_height*arena_tile_size-1 }
  if a.pac.Pos.Y >= a.Pos.Y+arena_height*arena_tile_size { a.pac.Pos.Y = a.Pos.Y+1 }

  // Update small coins
  var p_k utils.Vector2DInt64
  flag := false
  for k, c := range(a.sc) {
    if c.CheckPacmanInteraction(a.pac) {
      p_k = k
      flag = true
      break
    }
  }
  if flag {
    delete(a.sc, p_k)
    j, i := utils.Key2DToXY(p_k)
    a.Matrix[i][j] = 0
    a.Score += 10
  }

  // Update big coins
  flag = false
  for k, c := range(a.bc) {
    if c.CheckPacmanInteraction(a.pac) {
      p_k = k
      flag = true
      break
    }
  }
  if flag {
    delete(a.bc, p_k)
    j, i := utils.Key2DToXY(p_k)
    a.Matrix[i][j] = 0
    a.Score += 50

    a.ghost_red.Scare()
    a.ghost_blue.Scare()
    a.ghost_orange.Scare()
    a.ghost_pink.Scare()
  }

  // Update ghosts
  a.ghost_red.Update(a)
  a.ghost_blue.Update(a)
  a.ghost_orange.Update(a)
  a.ghost_pink.Update(a)

  // Check if they won
  if len(a.bc) == 0 && len(a.sc) == 0 {
    a.totalPlayTime += time.Since(a.lastPlayStart)
    a.state = ARENA_WIN_STATE
  }
}

func (a *Arena) Update(){
  if a.state == ARENA_PLAYING_STATE {a.GameUpdate()}
}

func (a *Arena) lazyBuildWall(i, j int) error {
  if _, ok := a.w[utils.Key2D(i, j)]; !ok {
    var err error
    a.w[utils.Key2D(i, j)], err = NewWall(&utils.Vector2DFloat64{X: float64(j*wall_spriteWidth)+a.Pos.X, Y: float64(i*wall_spriteHeight)+a.Pos.Y})
    if err != nil {
      return err
    }
  }

  return nil
}

func (a *Arena) lazyBuildBigCoin(i, j int) error {
  if _, ok := a.bc[utils.Key2D(i, j)]; !ok {
    var err error
    a.bc[utils.Key2D(i, j)], err = NewBigCoin(&utils.Vector2DFloat64{X: float64(j*big_coin_spriteWidth)+a.Pos.X, Y: float64(i*big_coin_spriteHeight)+a.Pos.Y})
    if err != nil {
      return err
    }
  }

  return nil
}

func (a *Arena) lazyBuildSmallCoin(i, j int) error {
  if _, ok := a.sc[utils.Key2D(i, j)]; !ok {
    var err error
    a.sc[utils.Key2D(i, j)], err = NewSmallCoin(&utils.Vector2DFloat64{X: float64(j*small_coin_spriteWidth)+a.Pos.X, Y: float64(i*small_coin_spriteHeight)+a.Pos.Y})
    if err != nil {
      return err
    }
  }
  return nil
}

func (a *Arena) Draw(screen *ebiten.Image) {
  for i := 0; i < arena_height; i++ {
    for j := 0; j < arena_width; j++ {
      if a.Matrix[i][j] == 1 {
        if _, ok := a.w[utils.Key2D(i, j)]; !ok {
          a.lazyBuildWall(i, j)
        }
        a.w[utils.Key2D(i, j)].Draw(screen)
        continue
      }
      if a.Matrix[i][j] == 2 {
        if _, ok := a.bc[utils.Key2D(i, j)]; !ok {
          a.lazyBuildBigCoin(i, j)
        }
        a.bc[utils.Key2D(i, j)].Draw(screen)
        continue
      }
      if a.Matrix[i][j] == 3 {
        if _, ok := a.sc[utils.Key2D(i, j)]; !ok {
          a.lazyBuildSmallCoin(i, j)
        }
        a.sc[utils.Key2D(i, j)].Draw(screen)
        continue
      }
    }
  }
  a.pac.Draw(screen)
  a.ghost_red.Draw(screen)
  a.ghost_blue.Draw(screen)
  a.ghost_pink.Draw(screen)
  a.ghost_orange.Draw(screen)

  op := &text.DrawOptions{}
  op.GeoM.Translate(arena_width*arena_tile_size, 0)
  op.ColorScale.ScaleWithColor(color.White)
  op.LineSpacing = ARENA_SCORE_FONT_SIZE
  op.PrimaryAlign = text.AlignEnd
  text.Draw(screen, fmt.Sprintf("SCORE: %d", a.Score), &text.GoTextFace{
    Source: a.arcadeFaceSource,
    Size: ARENA_SCORE_FONT_SIZE,
  }, op)

  op = &text.DrawOptions{}
  op.GeoM.Translate(0, 0)
  op.ColorScale.ScaleWithColor(color.White)
  op.LineSpacing = ARENA_SCORE_FONT_SIZE
  op.PrimaryAlign = text.AlignStart
  text.Draw(screen, fmt.Sprintf("LIVES: %d", a.pac.Lives), &text.GoTextFace{
    Source: a.arcadeFaceSource,
    Size: ARENA_SCORE_FONT_SIZE,
  }, op)

  if a.state == ARENA_RESET_STATE {
    a.start_button.Draw(screen)
  }

  if a.state == ARENA_START_STATE {
    a.start_text_box.Draw(screen)

    op = &text.DrawOptions{}
    op.GeoM.Translate(arena_width*arena_tile_size/2, 150)
    op.ColorScale.ScaleWithColor(color.White)
    op.LineSpacing = ARENA_SCORE_FONT_SIZE
    op.PrimaryAlign = text.AlignCenter
    text.Draw(screen, fmt.Sprintf("PACMAN V.42\n\nThis is a harder clone\nof the classic arcade\ngame Pacman\n\nEnemy movement is the\nsame as the classic\npacman however killed\nenemies follow the\noptimal path to their\nhome square. Enemies\nalso immediately leave\ntheir home square.\nFinally, enemy scatter\nmode is set on\na decay function\nbased on your score.\n\n\nPress (q) to exit"), &text.GoTextFace{
      Source: a.arcadeFaceSource,
      Size: ARENA_SCORE_FONT_SIZE,
    }, op)
  }

  if a.state == ARENA_END_STATE {
    a.end_text_box.Draw(screen)

    op = &text.DrawOptions{}
    op.GeoM.Translate(arena_width*arena_tile_size/2, 220)
    op.ColorScale.ScaleWithColor(color.White)
    op.LineSpacing = ARENA_SCORE_FONT_SIZE
    op.PrimaryAlign = text.AlignCenter
    text.Draw(screen, fmt.Sprintf("GAME OVER\n\n\nSCORE: %d\n\nTIME PLAYED: %.2f s\n\nSCORE PER SECOND: %.2f\n\n\nPress (q) to exit)", a.Score, a.totalPlayTime.Seconds(), float64(a.Score)/a.totalPlayTime.Seconds()), &text.GoTextFace{
      Source: a.arcadeFaceSource,
      Size: ARENA_SCORE_FONT_SIZE,
    }, op)
  }

  if a.state == ARENA_WIN_STATE {
    a.win_text_box.Draw(screen)

    op = &text.DrawOptions{}
    op.GeoM.Translate(arena_width*arena_tile_size/2, 200)
    op.ColorScale.ScaleWithColor(color.White)
    op.LineSpacing = ARENA_SCORE_FONT_SIZE
    op.PrimaryAlign = text.AlignCenter
    text.Draw(screen, fmt.Sprintf("WINNER\n\nWinner winner\nchicken dinner\n\nSCORE: %d\n\nTIME PLAYED: %.2f s\n\nSCORE PER SECOND: %.2f\n\n\nPress (q) to exit)", a.Score, a.totalPlayTime.Seconds(), float64(a.Score)/a.totalPlayTime.Seconds()), &text.GoTextFace{
      Source: a.arcadeFaceSource,
      Size: ARENA_SCORE_FONT_SIZE,
    }, op)
  }
}
