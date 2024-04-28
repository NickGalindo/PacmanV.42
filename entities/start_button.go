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
  start_button_assetWidth = 4990
  start_button_assetHeight = 237
  start_button_numTiles = 5
  start_button_tileWidth = 998
  start_button_tileHeight = 237
  start_button_spriteHeight = 47.49499
  start_button_spriteWidth = 200
)

type StartButton struct{
  Pos *utils.Vector2DFloat64
  frame int
  asset  *ebiten.Image
}

func NewStartButton(pos *utils.Vector2DFloat64) (*StartButton, error){
  file, err := os.Open(filepath.Join(utils.ASSET_DIR, "full_start_button.png"))
  if err != nil {
    return nil, err
  }
  defer file.Close()

  img, _, err := image.Decode(file)
  if err != nil{
    return nil, err
  }

  return &StartButton{
    Pos: pos,
    frame: 0,
    asset: ebiten.NewImageFromImage(img),
  }, nil
}

func (st_bt *StartButton) Draw(screen *ebiten.Image) {
  op := &ebiten.DrawImageOptions{}
  op.GeoM.Translate(-float64(start_button_tileWidth)/2, -float64(start_button_tileHeight)/2)
  op.GeoM.Scale(float64(start_button_spriteWidth)/float64(start_button_tileWidth), float64(start_button_spriteHeight)/float64(start_button_tileHeight))
  op.GeoM.Translate(st_bt.Pos.X, st_bt.Pos.Y)

  ax, ay := (st_bt.frame/6)*start_button_tileWidth, 0
  st_bt.frame += 1
  if (st_bt.frame/6) >= start_button_numTiles {
    st_bt.frame = -1
  }
  screen.DrawImage(st_bt.asset.SubImage(image.Rect(ax, ay, ax+start_button_tileWidth, ay+start_button_tileHeight)).(*ebiten.Image), op)
}
