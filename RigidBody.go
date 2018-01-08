package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

// ForceMode is the type of force to apply to a RigidBody
type ForceMode int8

const (
	// ConstantForce adds a continuous force, using mass
	ConstantForce ForceMode = iota
	// Acceleration adds a continuous force, ignoring mass
	Acceleration = iota
	// Impulse adds an instant force, using mass
	Impulse = iota
	// VelocityChange adds an instant force, ignoring mass
	VelocityChange = iota
)

// RigidBody is a physics body implemented with Rigid Body dynamics
type RigidBody struct {
	Parent       *Actor
	Radius       float32
	Mass         float32
	Velocity     mgl32.Vec3
	Acceleration mgl32.Vec3
}

// NewRigidBody creates a new RigidBody with appropriate defaults
func NewRigidBody() *RigidBody {
	return &RigidBody{
		Parent:       nil,
		Radius:       1.0,
		Mass:         1.0,
		Velocity:     mgl32.Vec3{0, 0, 0},
		Acceleration: mgl32.Vec3{0, 0, 0},
	}
}

// ApplyForce adds a force to the object, how it is added depenends on the mode
func (rb *RigidBody) ApplyForce(force mgl32.Vec3, mode ForceMode) {
	switch mode {
	case ConstantForce:
		force = force.Mul(1.0 / rb.Mass)
		rb.Acceleration = rb.Acceleration.Add(force)
	case Acceleration:
		rb.Acceleration = rb.Acceleration.Add(force)
	case Impulse:
		force = force.Mul(1.0 / rb.Mass)
		rb.Velocity = rb.Velocity.Add(force)
	case VelocityChange:
		rb.Velocity = rb.Velocity.Add(force)
	}
}

var LowerBound = mgl32.Vec3{0, 0, 0}
var UpperBound = mgl32.Vec3{100, 100, 100}

const Bounce = float32(0.5)
const Friction = float32(0.7)

func (rb *RigidBody) Update(delta, elapsed float32) {
	x, y, z := rb.Parent.Transform.Position.Add(rb.Velocity.Mul(delta)).Elem()
	vx, vy, vz := rb.Velocity.Add(rb.Acceleration.Mul(elapsed)).Elem()
	ax, ay, az := rb.Acceleration.Elem()

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
	pos := rb.Parent.Transform.Position
	otherPos := other.Parent.Transform.Position
	momentum := rb.Velocity.Add(other.Velocity).Len()

	rb.ApplyForce(pos.Add(otherPos.Mul(-1.0)).Normalize().Mul(momentum), Impulse)
	other.ApplyForce(otherPos.Add(pos.Mul(-1.0)).Normalize().Mul(momentum), Impulse)
}
