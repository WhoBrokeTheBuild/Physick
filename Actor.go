package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}

func NewTransform() Transform {
	return Transform{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.Vec3{0, 0, 0},
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

func (t Transform) GetMatrix() mgl32.Mat4 {
	return mgl32.Scale3D(t.Scale.X(), t.Scale.Y(), t.Scale.Z()).
		Mul4(mgl32.Rotate3DX(t.Rotation.X()).Mat4()).
		Mul4(mgl32.Rotate3DX(t.Rotation.Y()).Mat4()).
		Mul4(mgl32.Rotate3DX(t.Rotation.Z()).Mat4()).
		Mul4(mgl32.Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z()))
}

type Actor struct {
	Transform Transform
	RigidBody *RigidBody

	models []*Model
}

func NewActor() *Actor {
	actor := &Actor{
		Transform: NewTransform(),
		models:    []*Model{},
		RigidBody: NewRigidBody(),
	}
	actor.RigidBody.Parent = actor

	return actor
}

func (actor *Actor) AddModel(model *Model) {
	actor.models = append(actor.models, model)
}

func (actor *Actor) Update(delta, elapsed float32) {
	actor.RigidBody.Update(delta, elapsed)
}

func (actor *Actor) Render(shader *Shader) {
	shader.Use()

	model := actor.Transform.GetMatrix()
	gl.UniformMatrix4fv(shader.GetUniformLocation("u_Model"), 1, false, &model[0])

	for i := 0; i < len(actor.models); i++ {
		actor.models[i].Render()
	}
}
