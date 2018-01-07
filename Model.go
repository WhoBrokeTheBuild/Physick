package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Model struct {
	glVBO  uint32
	glVAO  uint32
	glType uint32
	count  int32
}

func NewCubeModel(shader *Shader, size float32) *Model {
	var vao uint32
	var vbo uint32

	cubeVertices := []float32{
		//  X, Y, Z, U, V
		// Bottom
		-size, -size, -size, 0.0, 0.0,
		size, -size, -size, 1.0, 0.0,
		-size, -size, size, 0.0, 1.0,
		size, -size, -size, 1.0, 0.0,
		size, -size, size, 1.0, 1.0,
		-size, -size, size, 0.0, 1.0,

		// Top
		-size, size, -size, 0.0, 0.0,
		-size, size, size, 0.0, 1.0,
		size, size, -size, 1.0, 0.0,
		size, size, -size, 1.0, 0.0,
		-size, size, size, 0.0, 1.0,
		size, size, size, 1.0, 1.0,

		// Front
		-size, -size, size, 1.0, 0.0,
		size, -size, size, 0.0, 0.0,
		-size, size, size, 1.0, 1.0,
		size, -size, size, 0.0, 0.0,
		size, size, size, 0.0, 1.0,
		-size, size, size, 1.0, 1.0,

		// Back
		-size, -size, -size, 0.0, 0.0,
		-size, size, -size, 0.0, 1.0,
		size, -size, -size, 1.0, 0.0,
		size, -size, -size, 1.0, 0.0,
		-size, size, -size, 0.0, 1.0,
		size, size, -size, 1.0, 1.0,

		// Left
		-size, -size, size, 0.0, 1.0,
		-size, size, -size, 1.0, 0.0,
		-size, -size, -size, 0.0, 0.0,
		-size, -size, size, 0.0, 1.0,
		-size, size, size, 1.0, 1.0,
		-size, size, -size, 1.0, 0.0,

		// Right
		size, -size, size, 1.0, 1.0,
		size, -size, -size, 1.0, 0.0,
		size, size, -size, 0.0, 0.0,
		size, -size, size, 1.0, 1.0,
		size, size, -size, 0.0, 0.0,
		size, size, size, 0.0, 1.0,
	}

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, gl.Ptr(cubeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(shader.GetAttribLocation("a_Vertex"))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	txcdAttrib := uint32(shader.GetAttribLocation("a_Texcoord"))
	gl.EnableVertexAttribArray(txcdAttrib)
	gl.VertexAttribPointer(txcdAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	model := &Model{
		glVAO:  vao,
		glVBO:  vbo,
		glType: gl.TRIANGLES,
		count:  6 * 2 * 3,
	}

	return model
}

func NewPlaneModel(shader *Shader, size float32) *Model {
	model := &Model{
		glVAO:  0,
		glVBO:  0,
		glType: gl.TRIANGLES,
		count:  0,
	}

	return model
}

func (model *Model) Render() {
	gl.BindVertexArray(model.glVAO)
	gl.DrawArrays(model.glType, 0, model.count)
}
