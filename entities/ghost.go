package entities

import (
	"image"
	_ "image/png"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/NickGalindo/Pacman/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
  ghost_death_assetWidth = 120
  ghost_death_assetHeight = 34
  ghost_death_num_tiles = 1
  ghost_death_tileWidth = 30
  ghost_death_tileHeight = 34

  ghost_scared_assetWidth = 120
  ghost_scared_assetHeight = 34
  ghost_scared_num_tiles = 2
  ghost_scared_tileWidth = 30
  ghost_scared_tileHeight = 34

  ghost_assetWidth = 240
  ghost_assetHeight = 34
  ghost_num_tiles = 2
  ghost_tileWidth = 30
  ghost_tileHeight = 34
  ghost_spriteHeight = 20
  ghost_spriteWidth = 20

  GHOST_DEATH_STATE = 0
  GHOST_SCARED_STATE = 1
  GHOST_ALIVE_STATE = 2
  
  GHOST_DIR_LEFT = 0
  GHOST_DIR_RIGHT = 2
  GHOST_DIR_UP = 4
  GHOST_DIR_DOWN = 6

  GHOST_SCATTER_MODE = false
  GHOST_ACTIVE_MODE = true
)


type Ghost struct{
  Pos *utils.Vector2DFloat64
  Collision_box *utils.Collider
  Interaction_box *utils.Collider
  asset  *ebiten.Image
  asset_scared *ebiten.Image
  asset_death *ebiten.Image
  frame int
  state uint8
  dir uint8
  scared_timer time.Time
  scared_total_time time.Duration
  mode bool
  pathing utils.PathingFunction
  path_dest utils.Vector2DFloat64
  speed float64
  mode_clock time.Time
  mode_active_time time.Duration
  mode_scatter_time time.Duration
  home utils.Vector2DFloat64
}

func (g *Ghost) GetColliderBox() (*utils.Collider) {
  return g.Collision_box
}
func (g *Ghost) GetInteractionBox() (*utils.Collider) {
  return g.Interaction_box
}

func NewGhost(pos *utils.Vector2DFloat64, ghost_sprite_path string, pathFunc utils.PathingFunction) (*Ghost, error) {
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, ghost_sprite_path))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  file_scared, err :=os.Open(filepath.Join(utils.ASSET_DIR, "ghost_scared.png"))
  if err != nil {
    return nil, err
  }
  defer file_scared.Close()

  img_scared, _, err := image.Decode(file_scared)
  if err != nil{
    return nil, err
  }

  file_death, err := os.Open(filepath.Join(utils.ASSET_DIR, "ghost_death.png"))
  if err != nil {
    return nil, err
  }
  defer file_death.Close()

  img_death, _, err := image.Decode(file_death)
  if err != nil{
    return nil, err
  }

  return &Ghost {
    Pos: pos,
    Collision_box: utils.NewCollider(pos, ghost_spriteWidth, ghost_spriteHeight),
    Interaction_box: utils.NewCollider(pos, ghost_spriteWidth, ghost_spriteHeight),
    asset: ebiten.NewImageFromImage(img),
    asset_scared: ebiten.NewImageFromImage(img_scared),
    asset_death: ebiten.NewImageFromImage(img_death),
    frame: 0,
    state: GHOST_ALIVE_STATE,
    dir: GHOST_DIR_UP,
    scared_timer: time.Now(),
    scared_total_time: 15*time.Second,
    pathing: pathFunc,
    path_dest: utils.Vector2DFloat64{X: 0, Y: 0},
    speed: 1,
    mode: GHOST_ACTIVE_MODE,
    mode_scatter_time: 2*time.Second,
    mode_active_time: 10*time.Second,
    home: *pos,
  }, nil
}

func (g *Ghost) Scare() {
  g.scared_timer = time.Now()
  g.state = GHOST_SCARED_STATE
}

func (g *Ghost) CalculatePathPoint(a *Arena) {
  cell_x, cell_y := a.FindCellFromXY(g.Pos.X, g.Pos.Y)
  home_x, home_y := a.FindCellFromXY(g.home.X, g.home.Y)
  pac_cell_x, pac_cell_y := a.FindCellFromXY(a.pac.Pos.X, a.pac.Pos.Y)
  red_ghost_cell_x, red_ghost_cell_y := a.FindCellFromXY(a.ghost_red.Pos.X, a.ghost_red.Pos.Y)

  var n_dir int
  var dest utils.Vector2DInt64
  if g.state == 0 {
    n_dir, dest = g.pathing(a.Matrix, cell_x, cell_y, int(g.dir)/2, home_x, home_y, int(a.pac.Dir), red_ghost_cell_x, red_ghost_cell_y, g.mode, int(g.state))

  } else {
    n_dir, dest = g.pathing(a.Matrix, cell_x, cell_y, int(g.dir)/2, pac_cell_x, pac_cell_y, int(a.pac.Dir), red_ghost_cell_x, red_ghost_cell_y, g.mode, int(g.state))
  }

  dest_x, dest_y := a.FindXYFromCell(int(dest.X), int(dest.Y))

  g.dir, g.path_dest = uint8(n_dir)*2, utils.Vector2DFloat64{X: dest_x, Y: dest_y}
}

func (g *Ghost) Update(a *Arena) {
  if g.state == GHOST_DEATH_STATE && *g.Pos == g.home {
    g.dir = GHOST_DIR_UP
    g.mode = GHOST_ACTIVE_MODE
    g.state = GHOST_ALIVE_STATE
    g.CalculatePathPoint(a)
  }
  if g.state == GHOST_SCARED_STATE && time.Since(g.scared_timer) > g.scared_total_time { g.state = GHOST_ALIVE_STATE }

  g.mode_scatter_time = time.Duration((5*math.Exp(-0.00258942004*float64(a.Score)))*1000000000)*time.Nanosecond

  if g.mode == GHOST_ACTIVE_MODE && time.Since(g.mode_clock) > g.mode_active_time {
    g.mode_clock = time.Now()
    g.mode = GHOST_SCATTER_MODE
    switch g.dir {
    case GHOST_DIR_UP:
      g.dir = GHOST_DIR_DOWN
    case GHOST_DIR_DOWN:
      g.dir = GHOST_DIR_UP
    case GHOST_DIR_LEFT:
      g.dir = GHOST_DIR_RIGHT
    case GHOST_DIR_RIGHT:
      g.dir = GHOST_DIR_LEFT
    }
    g.CalculatePathPoint(a)
  }
  if g.mode == GHOST_SCATTER_MODE && time.Since(g.mode_clock) > g.mode_scatter_time {
    g.mode_clock = time.Now()
    g.mode = GHOST_ACTIVE_MODE
    switch g.dir {
    case GHOST_DIR_UP:
      g.dir = GHOST_DIR_DOWN
    case GHOST_DIR_DOWN:
      g.dir = GHOST_DIR_UP
    case GHOST_DIR_LEFT:
      g.dir = GHOST_DIR_RIGHT
    case GHOST_DIR_RIGHT:
      g.dir = GHOST_DIR_LEFT
    }
    g.CalculatePathPoint(a)
  }

  if (g.dir == GHOST_DIR_DOWN || g.dir == GHOST_DIR_UP) && g.path_dest.Y == g.Pos.Y {g.CalculatePathPoint(a)}
  if (g.dir == GHOST_DIR_LEFT || g.dir == GHOST_DIR_RIGHT) && g.path_dest.X == g.Pos.X {g.CalculatePathPoint(a)}

  switch g.dir {
  case GHOST_DIR_UP:
    g.Pos.Y -= g.speed
  case GHOST_DIR_DOWN:
    g.Pos.Y += g.speed
  case GHOST_DIR_RIGHT:
    g.Pos.X += g.speed
  case GHOST_DIR_LEFT:
    g.Pos.X -= g.speed
  }


  if g.Pos.X < a.Pos.X { g.Pos.X = a.Pos.X+arena_width*arena_tile_size-1 }
  if g.Pos.X >= a.Pos.X+arena_width*arena_tile_size { g.Pos.X = a.Pos.X+1 }
  if g.Pos.Y < a.Pos.Y { g.Pos.Y = a.Pos.Y+arena_height*arena_tile_size-1 }
  if g.Pos.Y >= a.Pos.Y+arena_height*arena_tile_size { g.Pos.Y = a.Pos.Y+1 }
}

func (g *Ghost) Draw(screen *ebiten.Image) {
  var tileWidth, tileHeight, num_tiles int
  var ax, ay int
  var asset *ebiten.Image

  switch g.state {
  case GHOST_ALIVE_STATE:
    tileWidth, tileHeight, num_tiles = ghost_tileWidth, ghost_tileHeight, ghost_num_tiles
    asset = g.asset
    ax, ay = ((g.frame/6)+int(g.dir))*tileWidth, 0
  case GHOST_SCARED_STATE:
    tileWidth, tileHeight, num_tiles = ghost_scared_tileWidth, ghost_scared_tileHeight, ghost_scared_num_tiles
    asset = g.asset_scared

    if time.Since(g.scared_timer) > (g.scared_total_time*8)/10 {
      num_tiles = 2*ghost_scared_num_tiles
    }

    ax, ay = (g.frame/6)*tileWidth, 0
  case GHOST_DEATH_STATE:
    tileWidth, tileHeight, num_tiles = ghost_death_tileWidth, ghost_death_tileHeight, ghost_death_num_tiles
    asset = g.asset_death
    ax, ay = ((g.frame/6)+int(g.dir/2))*tileWidth, 0 
  }

  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(tileWidth)/2, -float64(tileHeight)/2)
  op.GeoM.Scale(float64(ghost_spriteWidth)/float64(tileWidth), float64(ghost_spriteHeight)/float64(tileHeight))
  op.GeoM.Translate(g.Pos.X, g.Pos.Y)


  g.frame += 1
  if (g.frame/6) >= num_tiles {
    g.frame = 0
  }
  screen.DrawImage(asset.SubImage(image.Rect(ax, ay, ax+tileWidth, ay+tileHeight)).(*ebiten.Image), op)
}
