package wogl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}

func NewTransform() Transform {
	var transform Transform
	transform.Position = mgl32.Vec3{0, 0, 0}
	transform.Rotation = mgl32.Vec3{0, 0, 0}
	transform.Scale = mgl32.Vec3{1, 1, 1}
	return transform
}
