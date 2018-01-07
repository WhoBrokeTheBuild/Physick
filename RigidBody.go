package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type RigidBody struct {
	Mass              float32
	Velocity          mgl32.Vec3
	Acceleration      mgl32.Vec3
	FixedAcceleration mgl32.Vec3
}

func NewRigidBody() *RigidBody {
	return &RigidBody{
		Mass:              1.0,
		Velocity:          mgl32.Vec3{0, 0, 0},
		Acceleration:      mgl32.Vec3{0, 0, 0},
		FixedAcceleration: mgl32.Vec3{0, 0, 0},
	}
}

func (rigidBody *RigidBody) ApplyForce(force mgl32.Vec3) {
	force = force.Mul(1.0 / rigidBody.Mass)
	rigidBody.Acceleration = rigidBody.Acceleration.Add(force).Add(rigidBody.FixedAcceleration)
}

func (rigidBody *RigidBody) ApplyConstantForce(force mgl32.Vec3) {
	rigidBody.FixedAcceleration = rigidBody.FixedAcceleration.Add(force)
	rigidBody.Acceleration = rigidBody.Acceleration.Add(rigidBody.FixedAcceleration)
}

func (rigidBody *RigidBody) Update(delta float32) {
	rigidBody.Velocity = rigidBody.Velocity.Add(rigidBody.Acceleration.Mul(delta))
}
