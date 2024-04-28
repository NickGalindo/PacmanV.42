package entities

import (
	"image"
	_ "image/png"
	"math"
	"os"

	"path/filepath"

	"github.com/NickGalindo/Pacman/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
  pacman_death_assetWidth = 352
  pacman_death_assetHeight = 32
  pacman_death_num_tiles = 11
  pacman_death_tileWidth = 32
  pacman_death_tileHeight = 32

  pacman_assetWidth = 128
  pacman_assetHeight = 32
  pacman_num_tiles = 4
  pacman_tileWidth = 32
  pacman_tileHeight = 32
  pacman_spriteHeight = 20
  pacman_spriteWidth = 20

  PACMAN_ALIVE_STATE = true
  PACMAN_DEAD_STATE = false

  PACMAN_DIR_UP = 0
  PACMAN_DIR_RIGHT = 1
  PACMAN_DIR_DOWN = 2
  PACMAN_DIR_LEFT = 3
)

type Pacman struct{
  Pos *utils.Vector2DFloat64
  Lives int64
  Collision_box *utils.Collider
  Interaction_box *utils.Collider
  asset  *ebiten.Image
  asset_death *ebiten.Image
  frame int
  state bool
  speed float64
  Dir uint8
  last_dir uint8
}

func (p *Pacman) GetColliderBox() (*utils.Collider) {
  return p.Collision_box
}
func (p *Pacman) GetInteractionBox() (*utils.Collider) {
  return p.Interaction_box
}

func NewPacman(pos *utils.Vector2DFloat64) (*Pacman, error) {
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, "pacman.png"))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  file_death, err := os.Open(filepath.Join(utils.ASSET_DIR, "pacman_death.png")) 
  if err != nil {
    return nil, err
  }
  defer file_death.Close()

  img_death, _, err := image.Decode(file_death)
  if err != nil{
    return nil, err
  }

  return &Pacman {
    Pos: pos,
    Lives: 3,
    Collision_box: utils.NewCollider(pos, pacman_spriteWidth, pacman_spriteHeight),
    Interaction_box: utils.NewCollider(pos, pacman_spriteWidth-10, pacman_spriteHeight-10),
    asset: ebiten.NewImageFromImage(img),
    asset_death: ebiten.NewImageFromImage(img_death),
    frame: 0,
    state: PACMAN_ALIVE_STATE,
    speed: 1,
    Dir: PACMAN_DIR_RIGHT,
    last_dir: PACMAN_DIR_RIGHT,
  }, nil
}

func (p *Pacman) CheckCollision(walls map[utils.Vector2DInt64]*Wall) bool {
  for k := range walls {
    if p.Collision_box.CheckCollision(walls[k].Collision_box) { return true }
  } 
  return false
}

func (p *Pacman) CheckInteraction(ghosts []*Ghost, score *int64) bool {
  for _, g := range ghosts {
    if p.Interaction_box.CheckCollision(g.Interaction_box) {
      if g.state == GHOST_DEATH_STATE {continue}
      if g.state == GHOST_SCARED_STATE {
        g.state = GHOST_DEATH_STATE
        *score = *score + 69
      } else {
        p.state = PACMAN_DEAD_STATE
        p.frame = 0
        p.Lives -= 1
        return true
      }
    }
  }
  return false
}

func (p *Pacman) HandlePressedKey(key ebiten.Key) {
  if p.state != PACMAN_ALIVE_STATE{return}

  switch key {
  case ebiten.KeyArrowUp:
    p.Dir = PACMAN_DIR_UP
  case ebiten.KeyArrowDown:
    p.Dir = PACMAN_DIR_DOWN
  case ebiten.KeyArrowRight:
    p.Dir = PACMAN_DIR_RIGHT
  case ebiten.KeyArrowLeft:
    p.Dir = PACMAN_DIR_LEFT
  }
}

func (p *Pacman) Update(walls map[utils.Vector2DInt64]*Wall, ghosts []*Ghost, score *int64) {
  for i := 0; i < 2; i++ {
    switch p.Dir {
    case PACMAN_DIR_UP:
      p.Pos.Y -= p.speed
    case PACMAN_DIR_DOWN:
      p.Pos.Y += p.speed
    case PACMAN_DIR_RIGHT:
      p.Pos.X += p.speed
    case PACMAN_DIR_LEFT:
      p.Pos.X -= p.speed
    }

    if p.CheckCollision(walls) {
      switch p.Dir {
      case PACMAN_DIR_UP:
        p.Pos.Y -= -p.speed
      case PACMAN_DIR_DOWN:
        p.Pos.Y += -p.speed
      case PACMAN_DIR_RIGHT:
        p.Pos.X += -p.speed
      case PACMAN_DIR_LEFT:
        p.Pos.X -= -p.speed
      }
      if p.Dir != p.last_dir {
        p.Dir = p.last_dir
      }
    } else {
      p.last_dir = p.Dir
      break
    }
  }

  p.CheckInteraction(ghosts, score)
}

func (p *Pacman) Draw(screen *ebiten.Image) {
  var tileWidth, tileHeight, num_tiles int
  var asset *ebiten.Image

  switch p.state {
  case PACMAN_ALIVE_STATE:
    tileWidth, tileHeight, num_tiles = pacman_tileWidth, pacman_tileHeight, pacman_num_tiles
    asset = p.asset
  case PACMAN_DEAD_STATE:
    tileWidth, tileHeight, num_tiles = pacman_death_tileWidth, pacman_death_tileHeight, pacman_death_num_tiles
    asset = p.asset_death
  }

  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(tileWidth)/2, -float64(tileHeight)/2)
  op.GeoM.Rotate(float64(p.Dir) * 1.0/2.0 * math.Pi)
  op.GeoM.Scale(float64(pacman_spriteWidth)/float64(tileWidth), float64(pacman_spriteHeight)/float64(tileHeight))
  op.GeoM.Translate(p.Pos.X, p.Pos.Y)

  ax, ay := (p.frame/6)*tileWidth, 0
  p.frame += 1
  if (p.frame/6) >= num_tiles {
    p.frame = -1
  }
  screen.DrawImage(asset.SubImage(image.Rect(ax, ay, ax+tileWidth, ay+tileHeight)).(*ebiten.Image), op)
}
