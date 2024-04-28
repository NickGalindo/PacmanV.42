package entities

import (
	"image"
	_ "image/png"
	"os"
  "path/filepath"
	"github.com/NickGalindo/Pacman/utils"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
  small_coin_assetWidth = 64
  small_coin_assetHeight = 16
  small_coin_num_tiles = 4
  small_coin_tileWidth = 16
  small_coin_tileHeight = 16
  small_coin_spriteHeight = 20
  small_coin_spriteWidth = 20
)

type SmallCoin struct{
  Pos *utils.Vector2DFloat64
  Interaction_box *utils.Collider
  asset  *ebiten.Image
  frame int
}

func (sc *SmallCoin) GetInteractionBox() (*utils.Collider) {
  return sc.Interaction_box
}

func NewSmallCoin(pos *utils.Vector2DFloat64) (*SmallCoin, error) {
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, "small_coin.png"))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  return &SmallCoin {
    Pos: pos,
    Interaction_box: utils.NewCollider(pos, small_coin_spriteWidth, small_coin_spriteHeight),
    asset: ebiten.NewImageFromImage(img),
    frame: 0,
  }, nil
}

func (sc *SmallCoin) CheckPacmanInteraction(p *Pacman) bool {
  return sc.Interaction_box.CheckCollision(p.Interaction_box)
}

func (sc *SmallCoin) Draw(screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(small_coin_tileWidth)/2, -float64(small_coin_tileHeight)/2)
  op.GeoM.Scale(float64(small_coin_spriteWidth)/float64(small_coin_tileWidth), float64(small_coin_spriteHeight)/float64(small_coin_tileHeight))
  op.GeoM.Translate(sc.Pos.X, sc.Pos.Y)

  ax, ay := (sc.frame/12)*small_coin_tileWidth, 0
  sc.frame += 1
  if (sc.frame/12) >= small_coin_num_tiles {
    sc.frame = 0
  }

  screen.DrawImage(sc.asset.SubImage(image.Rect(ax, ay, ax+small_coin_tileWidth, ay+small_coin_tileHeight)).(*ebiten.Image), op)
}
