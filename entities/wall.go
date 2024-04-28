package entities

import (
  "image"
  _ "image/png"
  "os"
  "path/filepath"
  "math/rand/v2"
	"github.com/hajimehoshi/ebiten/v2"
  "github.com/NickGalindo/Pacman/utils"
)

const (
  wall_assetWidth = 72
  wall_assetHeight = 18
  wall_numTiles = 4
  wall_tileWidth = 18
  wall_tileHeight = 18
  wall_spriteHeight = 20
  wall_spriteWidth = 20
)

type Wall struct{
  Pos *utils.Vector2DFloat64
  Collision_box *utils.Collider
  asset  *ebiten.Image
}

func (w *Wall) GetColliderBox() (*utils.Collider) {
  return w.Collision_box
}

func NewWall(pos *utils.Vector2DFloat64) (*Wall, error){
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, "walls2.png"))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  i := rand.IntN(wall_numTiles)
  ax, ay := i*wall_tileWidth, 0

  return &Wall {
    Pos: pos,
    Collision_box: utils.NewCollider(pos, wall_spriteWidth, wall_spriteHeight),
    asset: ebiten.NewImageFromImage(img).SubImage(image.Rect(ax, ay, ax+wall_tileWidth, ay+wall_tileHeight)).(*ebiten.Image),
  }, nil
}

func (w *Wall) Draw(screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(wall_tileWidth)/2, -float64(wall_tileHeight)/2)
  op.GeoM.Scale(float64(wall_spriteWidth)/float64(wall_tileWidth), float64(wall_spriteHeight)/float64(wall_tileHeight))
  op.GeoM.Translate(w.Pos.X, w.Pos.Y)

  screen.DrawImage(w.asset, op)
}
