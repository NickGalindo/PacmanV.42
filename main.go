package main

import (
	"log"

	"image/color"

	"github.com/NickGalindo/Pacman/entities"
	"github.com/NickGalindo/Pacman/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var a *entities.Arena

type Game struct{
  keys []ebiten.Key
}

func (g *Game) handlePressedKeys() {
  g.keys = inpututil.AppendPressedKeys(g.keys[:0])

  if len(g.keys) == 0 {
    return
  }

  a.HandlePressedKey(g.keys[0])

}

func (g *Game) Update() error {
  g.handlePressedKeys()

  a.Update()

  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
  screen.Fill(color.RGBA{4, 16, 54, 0xff})
  //ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f", ebiten.ActualFPS()))
  a.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  return 460, 600
}

func main() {
  var err error
  a, err = entities.NewArena(&utils.Vector2DFloat64{X: 10, Y: 30})
  if err != nil {
    log.Fatal(err)
  }

  a.Init()

  ebiten.SetWindowSize(460, 600)
  ebiten.SetWindowTitle("Helol Worldo")
  err = ebiten.RunGame(&Game{})
  if err != nil{
    log.Fatal(err)
  }
}
