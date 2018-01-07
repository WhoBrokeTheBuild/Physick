package main

import (
	"log"
	"math"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

	cubeData := []float32{
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
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeData)*4, gl.Ptr(cubeData), gl.STATIC_DRAW)

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

// TODO Make work
func NewSphereModel(shader *Shader, radius float32) *Model {
	const rings = 5
	const sectors = 10

	const count = rings * sectors * 6

	tmpVerts := make([]mgl32.Vec3, (rings+1)*(sectors+1))
	tmpTxcds := make([]mgl32.Vec2, (rings+1)*(sectors+1))

	data := make([]float32, count*(3+2))

	R := 1.0 / float64(rings-1.0)
	S := 1.0 / float64(sectors-1.0)

	i := 0
	for r := 0; r < rings+1; r++ {
		for s := 0; s < sectors+1; s++ {
			x := math.Sin((-math.Pi * 0.5) + (math.Pi * float64(r) * R))
			y := math.Cos(2.0*math.Pi*float64(s)*S) * math.Sin(math.Pi*float64(r)*R)
			z := math.Cos(2.0*math.Pi*float64(s)*S) * math.Sin(math.Pi*float64(r)*R)

			tmpVerts[i] = mgl32.Vec3{
				radius * float32(x),
				radius * float32(y),
				radius * float32(z),
			}
			//tmpNorms[i] = mgl32.Vec3{
			//    float32(x),
			//    float32(y),
			//    float32(z),
			//}
			tmpTxcds[i] = mgl32.Vec2{
				float32(s) * float32(S),
				float32(r) * float32(R),
			}
			i++
		}
	}

	i = 0
	for r := 0; r < rings; r++ {
		for s := 0; s < sectors; s++ {
			log.Println(r, s, count, i)
			data[i+0] = tmpVerts[r*sectors+s][0]
			data[i+1] = tmpVerts[r*sectors+s][1]
			data[i+2] = tmpVerts[r*sectors+s][2]
			data[i+3] = tmpTxcds[r*sectors+s][0]
			data[i+4] = tmpTxcds[r*sectors+s][1]
			i += 5
			data[i+0] = tmpVerts[r*sectors+(s+1)][0]
			data[i+1] = tmpVerts[r*sectors+(s+1)][1]
			data[i+2] = tmpVerts[r*sectors+(s+1)][2]
			data[i+3] = tmpTxcds[r*sectors+(s+1)][0]
			data[i+4] = tmpTxcds[r*sectors+(s+1)][1]
			i += 5
			data[i+0] = tmpVerts[(r+1)*sectors+s][0]
			data[i+1] = tmpVerts[(r+1)*sectors+s][1]
			data[i+2] = tmpVerts[(r+1)*sectors+s][2]
			data[i+3] = tmpTxcds[(r+1)*sectors+s][0]
			data[i+4] = tmpTxcds[(r+1)*sectors+s][1]
			i += 5
			data[i+0] = tmpVerts[r*sectors+(s+1)][0]
			data[i+1] = tmpVerts[r*sectors+(s+1)][1]
			data[i+2] = tmpVerts[r*sectors+(s+1)][2]
			data[i+3] = tmpTxcds[r*sectors+(s+1)][0]
			data[i+4] = tmpTxcds[r*sectors+(s+1)][1]
			i += 5
			data[i+0] = tmpVerts[(r+1)*sectors+s][0]
			data[i+1] = tmpVerts[(r+1)*sectors+s][1]
			data[i+2] = tmpVerts[(r+1)*sectors+s][2]
			data[i+3] = tmpTxcds[(r+1)*sectors+s][0]
			data[i+4] = tmpTxcds[(r+1)*sectors+s][1]
			i += 5
			data[i+0] = tmpVerts[(r+1)*sectors+(s+1)][0]
			data[i+1] = tmpVerts[(r+1)*sectors+(s+1)][1]
			data[i+2] = tmpVerts[(r+1)*sectors+(s+1)][2]
			data[i+3] = tmpTxcds[(r+1)*sectors+(s+1)][0]
			data[i+4] = tmpTxcds[(r+1)*sectors+(s+1)][1]
			i += 5
		}
	}

	var vao uint32
	var vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	vertAttrib := uint32(shader.GetAttribLocation("a_Vertex"))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 4, gl.PtrOffset(0))

	txcdAttrib := uint32(shader.GetAttribLocation("a_Texcoord"))
	gl.EnableVertexAttribArray(txcdAttrib)
	gl.VertexAttribPointer(txcdAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	model := &Model{
		glVAO:  vao,
		glVBO:  vbo,
		glType: gl.TRIANGLES,
		count:  count,
	}

	return model
}

func NewPlaneModel(shader *Shader, size float32) *Model {
	const rows = 10
	const cols = 10

	const count = (rows * cols) + (rows-1)*(cols-2)

	sqWidth := size / cols
	sqHeight := size / rows

	tmpVerts := make([]mgl32.Vec3, rows*cols)
	tmpTxcds := make([]mgl32.Vec2, rows*cols)

	data := make([]float32, count*(3+2))

	i := 0
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			tmpVerts[i] = mgl32.Vec3{float32(col) * sqWidth, 0.0, float32(row) * sqHeight}
			tmpTxcds[i] = mgl32.Vec2{float32(col) / float32(cols), float32(row) / float32(rows)}
			i++
		}
	}

	i = 0
	for row := 0; row < rows-1; row++ {
		if row&1 == 0 {
			// Even Rows
			for col := 0; col < cols; col++ {
				data[i+0] = tmpVerts[col+row*cols][0]
				data[i+1] = tmpVerts[col+row*cols][1]
				data[i+2] = tmpVerts[col+row*cols][2]
				data[i+3] = tmpTxcds[col+row*cols][0]
				data[i+4] = tmpTxcds[col+row*cols][1]
				i += 5
				data[i+0] = tmpVerts[col+(row+1)*cols][0]
				data[i+1] = tmpVerts[col+(row+1)*cols][1]
				data[i+2] = tmpVerts[col+(row+1)*cols][2]
				data[i+3] = tmpTxcds[col+(row+1)*cols][0]
				data[i+4] = tmpTxcds[col+(row+1)*cols][1]
				i += 5
			}
		} else {
			// Odd Rows
			for col := cols - 1; col > 0; col-- {
				data[i+0] = tmpVerts[col+(row+1)*cols][0]
				data[i+1] = tmpVerts[col+(row+1)*cols][1]
				data[i+2] = tmpVerts[col+(row+1)*cols][2]
				data[i+3] = tmpTxcds[col+(row+1)*cols][0]
				data[i+4] = tmpTxcds[col+(row+1)*cols][1]
				i += 5
				data[i+0] = tmpVerts[col-1+row*cols][0]
				data[i+1] = tmpVerts[col-1+row*cols][1]
				data[i+2] = tmpVerts[col-1+row*cols][2]
				data[i+3] = tmpTxcds[col-1+row*cols][0]
				data[i+4] = tmpTxcds[col-1+row*cols][1]
				i += 5
			}
		}
	}
	if rows&1 == 1 && rows > 2 {
		data[i+0] = tmpVerts[(rows-1)*cols][0]
		data[i+1] = tmpVerts[(rows-1)*cols][1]
		data[i+2] = tmpVerts[(rows-1)*cols][2]
		data[i+3] = tmpTxcds[(rows-1)*cols][0]
		data[i+4] = tmpTxcds[(rows-1)*cols][1]
		i += 5
	}

	var vao uint32
	var vbo uint32

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)

	vertAttrib := uint32(shader.GetAttribLocation("a_Vertex"))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	txcdAttrib := uint32(shader.GetAttribLocation("a_Texcoord"))
	gl.EnableVertexAttribArray(txcdAttrib)
	gl.VertexAttribPointer(txcdAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	model := &Model{
		glVAO:  vao,
		glVBO:  vbo,
		glType: gl.TRIANGLE_STRIP,
		count:  count,
	}

	return model
}

func (model *Model) Render() {
	gl.BindVertexArray(model.glVAO)
	gl.DrawArrays(model.glType, 0, model.count)
}
