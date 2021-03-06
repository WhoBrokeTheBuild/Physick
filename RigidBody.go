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

type Collider interface {
}

type SphereCollider struct {
	Radius float32
}

type BoxCollider struct {
	Size mgl32.Vec3
}

// RigidBody is a physics body implemented with Rigid Body dynamics
type RigidBody struct {
	Parent       *Actor
	Collider     Collider
	Mass         float32
	Velocity     mgl32.Vec3
	Acceleration mgl32.Vec3
}

// NewRigidBody creates a new RigidBody with appropriate defaults
func NewRigidBody() *RigidBody {
	return &RigidBody{
		Parent:       nil,
		Mass:         1.0,
		Velocity:     mgl32.Vec3{0, 0, 0},
		Acceleration: mgl32.Vec3{0, 0, 0},
	}
}

// Cleanup frees up resources
func (rb *RigidBody) Cleanup() {
	rb.Parent = nil
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

/*
var LowerBound = mgl32.Vec3{0, 0, 0}
var UpperBound = mgl32.Vec3{100, 100, 100}

const Bounce = float32(0.5)
const Friction = float32(0.7)
*/

func (rb *RigidBody) Update(delta, elapsed float32) {
	x, y, z := rb.Parent.Transform.Position.Add(rb.Velocity.Mul(delta)).Elem()
	vx, vy, vz := rb.Velocity.Add(rb.Acceleration.Mul(elapsed)).Elem()
	ax, ay, az := rb.Acceleration.Elem()

	/*
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
	*/

	//log.Println(vx, vy, vz)

	rb.Parent.Transform.Position = mgl32.Vec3{x, y, z}
	rb.Velocity = mgl32.Vec3{vx, vy, vz}
	rb.Acceleration = mgl32.Vec3{ax, ay, az}

	// Min/Max Velocity
	if rb.Velocity.Len() > 100.0 {
		rb.Velocity = rb.Velocity.Normalize().Mul(100.0)
	}
	if rb.Velocity.Len() < 0.001 {
		rb.Velocity = mgl32.Vec3{0, 0, 0}
	}
}

func (rb *RigidBody) CheckCollide(other *RigidBody) {
	pos := rb.Parent.Transform.Position
	otherPos := other.Parent.Transform.Position

	switch rb.Collider.(type) {
	case SphereCollider:
		switch other.Collider.(type) {
		case SphereCollider:
			dist := DistanceSquared(pos, otherPos)
			col := rb.Collider.(SphereCollider)
			otherCol := other.Collider.(SphereCollider)
			if dist < (col.Radius+otherCol.Radius)*(col.Radius+otherCol.Radius) {
				rb.Collide(other)
			}
		}
	case BoxCollider:

	}
}

func (rb *RigidBody) Collide(other *RigidBody) {
	diff := rb.Parent.Transform.Position.Add(other.Parent.Transform.Position.Mul(-1.0))

	x := diff.Normalize()
	v1 := rb.Velocity
	x1 := x.Dot(v1)
	v1x := x.Mul(x1)
	v1y := v1.Add(v1x.Mul(-1.0))
	m1 := rb.Mass

	x = x.Mul(-1.0)
	v2 := other.Velocity
	x2 := x.Dot(v2)
	v2x := x.Mul(x2)
	v2y := v2.Add(v2x.Mul(-1.0))
	m2 := other.Mass

	rb.Velocity = v1x.Mul((m1 - m2) / (m1 + m2)).Add(v2x.Mul((2 * m2) / (m1 + m2)).Add(v1y))
	other.Velocity = v1x.Mul((2 * m1) / (m1 + m2)).Add(v2x.Mul((m2 - m1) / (m1 + m2)).Add(v2y))
}
