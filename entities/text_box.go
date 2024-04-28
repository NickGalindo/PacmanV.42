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
  text_box_assetWidth = 460
  text_box_assetHeight = 460
  text_box_tileWidth = 460
  text_box_tileHeight = 460
  text_box_spriteHeight = 430
  text_box_spriteWidth = 430
)

type TextBox struct{
  Pos *utils.Vector2DFloat64
  asset  *ebiten.Image
}

func NewTextBox(pos *utils.Vector2DFloat64, asset_name string) (*TextBox, error){
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, asset_name))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  return &TextBox{
    Pos: pos,
    asset: ebiten.NewImageFromImage(img),
  }, nil
}

func (tb *TextBox) Draw(screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(text_box_tileWidth)/2, -float64(text_box_tileHeight)/2)
  op.GeoM.Scale(float64(text_box_spriteWidth)/float64(text_box_tileWidth), float64(text_box_spriteHeight)/float64(text_box_tileHeight))
  op.GeoM.Translate(tb.Pos.X, tb.Pos.Y)

  screen.DrawImage(tb.asset, op)
}
