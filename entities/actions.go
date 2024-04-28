package entities

import "github.com/NickGalindo/Pacman/utils"

type ColliderBoxed interface {
  GetColliderBox() *utils.Collider
}

type InteractionBoxed interface {
  GetInteractionBox() *utils.Collider
}
