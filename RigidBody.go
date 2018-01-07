package main

import "github.com/go-gl/mathgl/mgl32"

type Collider interface {
}

type BoxCollider struct {
	Position, Size mgl32.Vec3
}

type SphereCollider struct {
	Radius float32
}

type RigidBody struct {
	Fixed bool

	colliders []*Collider
}

func NewRigidBody() *RigidBody {
	return &RigidBody{
		Fixed: false,
	}
}

func (rigidBody *RigidBody) AddCollider(collider *Collider) {
	rigidBody.colliders = append(rigidBody.colliders, collider)
}

func (rigidBody *RigidBody) Update(delta float32) {

}
