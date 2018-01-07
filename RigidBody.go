package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type RigidBody struct {
	Parent            *Actor
	Radius            float32
	Mass              float32
	Velocity          mgl32.Vec3
	Acceleration      mgl32.Vec3
	FixedAcceleration mgl32.Vec3
}

func NewRigidBody() *RigidBody {
	return &RigidBody{
		Parent:            nil,
		Radius:            1.0,
		Mass:              1.0,
		Velocity:          mgl32.Vec3{0, 0, 0},
		Acceleration:      mgl32.Vec3{0, 0, 0},
		FixedAcceleration: mgl32.Vec3{0, 0, 0},
	}
}

func (rb *RigidBody) ApplyForce(force mgl32.Vec3) {
	force = force.Mul(1.0 / rb.Mass)
	rb.Acceleration = force.Add(rb.FixedAcceleration)
}

func (rb *RigidBody) ApplyConstantForce(force mgl32.Vec3) {
	rb.FixedAcceleration = rb.FixedAcceleration.Add(force)
}

var LowerBound = mgl32.Vec3{0, 0, 0}
var UpperBound = mgl32.Vec3{100, 100, 100}

const Bounce = float32(0.5)
const Friction = float32(0.7)

func (rb *RigidBody) Update(delta, elapsed float32) {
	x, y, z := rb.Parent.Transform.Position.Add(rb.Velocity.Mul(delta)).Elem()
	vx, vy, vz := rb.Velocity.Add(rb.Acceleration.Mul(delta)).Elem()
	ax, ay, az := rb.FixedAcceleration.Mul(delta).Elem()

	if x < LowerBound.X() {
		x = LowerBound.X()
		vx = -vx * Bounce
		//ax = -ax
	} else if x > UpperBound.X() {
		x = UpperBound.X()
		vx = -vx * Bounce
		//ax = -ax
	}

	if y < LowerBound.Y() {
		y = LowerBound.Y()
		vy = -vy * Bounce
		//ay = 0.0

		// Friction
		vx *= Friction
		vz *= Friction
	} else if y > UpperBound.Y() {
		y = UpperBound.Y()
		vy = -vy * Bounce
		//ay = -ay
	}

	if z < LowerBound.Z() {
		z = LowerBound.Z()
		vz = -vz * Bounce
		//az = -az
	} else if z > UpperBound.Z() {
		z = UpperBound.Z()
		vz = -vz * Bounce
		//az = -az
	}

	// Clamp to floor
	if vy < 0.0 && vy > 0.001 {
		vy = 0.0
	}

	//log.Println(vx, vy, vz)

	rb.Parent.Transform.Position = mgl32.Vec3{x, y, z}
	rb.Velocity = mgl32.Vec3{vx, vy, vz}
	rb.Acceleration = mgl32.Vec3{ax, ay, az}
}

func (rb *RigidBody) CheckCollide(other *RigidBody) {
	pos := rb.Parent.Transform.Position
	otherPos := other.Parent.Transform.Position

	dist := DistanceSquared(pos, otherPos)
	if dist < rb.Radius+other.Radius {
		rb.Collide(other)
	}
}

func (rb *RigidBody) Collide(other *RigidBody) {

}
