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

var LowerBound = mgl32.Vec3{0, 0, 0}
var UpperBound = mgl32.Vec3{100, 100, 100}

func NewActor() *Actor {
	actor := &Actor{
		Transform: NewTransform(),
		models:    []*Model{},
		RigidBody: NewRigidBody(),
	}

	return actor
}

func (actor *Actor) AddModel(model *Model) {
	actor.models = append(actor.models, model)
}

func (actor *Actor) Update(delta float32) {
	actor.RigidBody.Update(delta)

	bounce := float32(0.5)
	friction := float32(0.7)

	x, y, z := actor.Transform.Position.Add(actor.RigidBody.Velocity).Elem()
	vx, vy, vz := actor.RigidBody.Velocity.Elem()
	ax, ay, az := actor.RigidBody.Acceleration.Elem()

	if x < LowerBound.X() {
		x = LowerBound.X()
		vx = -vx * bounce
		//ax = -ax
	} else if x > UpperBound.X() {
		x = UpperBound.X()
		vx = -vx * bounce
		//ax = -ax
	}

	if y < LowerBound.Y() {
		y = LowerBound.Y()
		vy = -vy * bounce
		//ay = 0.0

		// Friction
		vx *= friction
		vz *= friction
	} else if y > UpperBound.Y() {
		y = UpperBound.Y()
		vy = -vy * bounce
		//ay = -ay
	}

	if z < LowerBound.Z() {
		z = LowerBound.Z()
		vz = -vz * bounce
		//az = -az
	} else if z > UpperBound.Z() {
		z = UpperBound.Z()
		vz = -vz * bounce
		//az = -az
	}

	// Clamp to floor
	if vy < 0.0 && vy > 0.001 {
		vy = 0.0
	}

	//log.Println(vx, vy, vz)

	actor.Transform.Position = mgl32.Vec3{x, y, z}
	actor.RigidBody.Velocity = mgl32.Vec3{vx, vy, vz}
	actor.RigidBody.Acceleration = mgl32.Vec3{ax, ay, az}
}

func (actor *Actor) Render(shader *Shader) {
	shader.Use()

	model := actor.Transform.GetMatrix()
	gl.UniformMatrix4fv(shader.GetUniformLocation("u_Model"), 1, false, &model[0])

	for i := 0; i < len(actor.models); i++ {
		actor.models[i].Render()
	}
}
