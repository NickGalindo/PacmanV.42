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
  big_coin_assetWidth = 80
  big_coin_assetHeight = 16
  big_coin_numTiles = 5
  big_coin_tileWidth = 16
  big_coin_tileHeight = 16
  big_coin_spriteHeight = 20
  big_coin_spriteWidth = 20
)

type BigCoin struct{
  Pos *utils.Vector2DFloat64
  Interaction_box *utils.Collider
  asset  *ebiten.Image
  frame int
}

func (bg *BigCoin) GetInteractionBox() (*utils.Collider) {
  return bg.Interaction_box
}

func NewBigCoin(pos *utils.Vector2DFloat64) (*BigCoin, error) {
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, "big_coin.png"))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  return &BigCoin {
    Pos: pos,
    Interaction_box: utils.NewCollider(pos, big_coin_spriteWidth, big_coin_spriteHeight),
    asset: ebiten.NewImageFromImage(img),
    frame: 0,
  }, nil
}

func (bc *BigCoin) CheckPacmanInteraction(p *Pacman) bool {
  return bc.Interaction_box.CheckCollision(p.Interaction_box)
}

func (bc *BigCoin) Draw(screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(big_coin_tileWidth)/2, -float64(big_coin_tileHeight)/2)
  op.GeoM.Scale(float64(big_coin_spriteWidth)/float64(big_coin_tileWidth), float64(big_coin_spriteHeight)/float64(big_coin_tileHeight))
  op.GeoM.Translate(bc.Pos.X, bc.Pos.Y)

  ax, ay := (bc.frame/12)*big_coin_tileWidth, 0
  bc.frame += 1
  if (bc.frame/12) >= big_coin_numTiles {
    bc.frame = 0
  }

  screen.DrawImage(bc.asset.SubImage(image.Rect(ax, ay, ax+big_coin_tileWidth, ay+big_coin_tileHeight)).(*ebiten.Image), op)
}
