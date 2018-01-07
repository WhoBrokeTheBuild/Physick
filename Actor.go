package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Actor struct {
	Transform mgl32.Mat4
	RigidBody RigidBody

	models []*Model
}

func NewActor() *Actor {
	actor := &Actor{
		Transform: mgl32.Ident4(),
		models:    []*Model{},
	}

	return actor
}

func (actor *Actor) AddModel(model *Model) {
	actor.models = append(actor.models, model)
}

func (actor *Actor) Render(shader *Shader) {
	shader.Use()
	gl.UniformMatrix4fv(shader.GetUniformLocation("u_Model"), 1, false, &actor.Transform[0])
	for i := 0; i < len(actor.models); i++ {
		actor.models[i].Render()
	}
}
