package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	VERT_ATTRIB = 0
	NORM_ATTRIB = 1
	TXCD_ATTRIB = 2
)

type modelGroup struct {
	DrawMode uint32
	Start    int32
	Count    int32
}

type Model struct {
	Transform mgl32.Mat4

	glVao  uint32
	glVbos [3]uint32
	groups []modelGroup
}

func NewModel() (*Model, error) {
	return &Model{
		Transform: mgl32.Ident4(),
		glVao:     0,
		glVbos:    [3]uint32{0, 0, 0},
	}, nil
}

func NewModelFromFile(filename string) (*Model, error) {
	model, err := NewModel()
	if err != nil {
		return model, err
	}
	err = model.LoadFromFile(filename)
	if err != nil {
		return model, err
	}
	return model, nil
}

func (model *Model) Cleanup() {
	gl.DeleteBuffers(3, &model.glVbos[0])
	gl.DeleteVertexArrays(1, &model.glVao)
}

func (model *Model) LoadFromFile(filename string) error {
	// Holds a material
	type MatDef struct {
		Ambient     mgl32.Vec3
		Diffuse     mgl32.Vec3
		Specular    mgl32.Vec3
		Shininess   float32
		Dissolve    float32
		AmbientMap  string
		SpecularMap string
		DiffuseMap  string
		BumpMap     string
	}

	// Holds a single face
	type Face struct {
		VertInds [3]int
		NormInds [3]int
		TxcdInds [3]int
	}

	// Holds a group of faces and a material
	type Group struct {
		Name     string
		Material string
		Faces    []Face
	}

	// Open the .obj file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Create a reader with a specific buffer size, needed by reader.ReadLine()
	reader := bufio.NewReader(bytes.NewReader(data))

	// Create a list of groups, and get a pointer to the first
	groups := []Group{{}}
	group := &groups[0]

	// Create the list of all Vertices, Normals, and Texture Coordinates
	allVerts := []mgl32.Vec3{}
	allNorms := []mgl32.Vec3{}
	allTxcds := []mgl32.Vec2{}

	var line string
	var count int

	tmpVec3 := mgl32.Vec3{}
	tmpVec2 := mgl32.Vec2{}
	tmpFace := Face{}

	tmp, _, err := reader.ReadLine()
	for ; err == nil; tmp, _, err = reader.ReadLine() {
		line = string(tmp)

		// Ignore empty lines and comments
		if len(line) == 0 || line[0] == '#' {
			continue
		}

		// Split on the first ' ', ignore half lines
		parts := strings.SplitN(line, " ", 2)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "usemtl":

			group.Material = parts[1]

		case "mtllib":

			if err != nil {
				return err
			}

		case "o":
			fallthrough
		case "g":

			if group.Name == "" {
				group.Name = parts[1]
			} else {
				groups = append(groups, Group{
					Name: parts[1],
				})
				group = &groups[len(groups)-1]
			}

		case "f":

			// Test for and parse faces in the 'v//vn v//vn v//vn' format
			if strings.Contains(parts[1], "//") {
				count, err = fmt.Sscanf(parts[1],
					"%d//%d %d//%d %d//%d",
					&tmpFace.VertInds[0],
					&tmpFace.NormInds[0],
					&tmpFace.VertInds[1],
					&tmpFace.NormInds[1],
					&tmpFace.VertInds[2],
					&tmpFace.NormInds[2],
				)
				if err != nil || count != 6 {
					return fmt.Errorf("Malformed OBJ file '%v' %v", line)
				}
				// Test for and parse faces in the 'v/vt v/vt v/vt' format
			} else if strings.Count(parts[1], "/") == 3 {
				count, err = fmt.Sscanf(parts[1],
					"%d/%d %d/%d %d/%d",
					&tmpFace.VertInds[0],
					&tmpFace.TxcdInds[0],
					&tmpFace.VertInds[1],
					&tmpFace.TxcdInds[1],
					&tmpFace.VertInds[2],
					&tmpFace.TxcdInds[2],
				)
				if err != nil || count != 6 {
					return fmt.Errorf("Malformed OBJ file '%v'", line)
				}
				// Test for and parse faces in the 'v/vt/vn v/vt/vn v/vt/vn' format
			} else if strings.Count(parts[1], "/") == 6 {
				count, err = fmt.Sscanf(parts[1],
					"%d/%d/%d %d/%d/%d %d/%d/%d",
					&tmpFace.VertInds[0],
					&tmpFace.TxcdInds[0],
					&tmpFace.NormInds[0],
					&tmpFace.VertInds[1],
					&tmpFace.TxcdInds[1],
					&tmpFace.NormInds[1],
					&tmpFace.VertInds[2],
					&tmpFace.TxcdInds[2],
					&tmpFace.NormInds[2],
				)
				if err != nil || count != 9 {
					return fmt.Errorf("Malformed OBJ file '%v'", line)
				}
			} else {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			group.Faces = append(group.Faces, tmpFace)

		case "v":

			count, err = fmt.Sscanf(parts[1], "%f %f %f", &tmpVec3[0], &tmpVec3[1], &tmpVec3[2])
			if err != nil || count != 3 {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			allVerts = append(allVerts, tmpVec3)

		case "vn":

			count, err = fmt.Sscanf(parts[1], "%f %f %f", &tmpVec3[0], &tmpVec3[1], &tmpVec3[2])
			if err != nil || count != 3 {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			allNorms = append(allNorms, tmpVec3)

		case "vt":

			count, err = fmt.Sscanf(parts[1], "%f %f", &tmpVec2[0], &tmpVec2[1])
			if err != nil || count != 2 {
				return fmt.Errorf("Malformed OBJ file '%v'", line)
			}

			allTxcds = append(allTxcds, tmpVec2)

		}
	}

	if group.Name == "" {
		group.Name = "default"
	}

	start := int32(0)
	verts := []float32{}
	norms := []float32{}
	txcds := []float32{}

	for g := range groups {
		group := &groups[g]
		for f := range group.Faces {
			face := &group.Faces[f]
			for i := 0; i < 3; i++ {
				// Adjust for negative indices
				if face.VertInds[i] < 0 {
					face.VertInds[i] += len(allVerts)
				}
				if face.NormInds[i] < 0 {
					face.NormInds[i] += len(allNorms)
				}
				if face.TxcdInds[i] < 0 {
					face.TxcdInds[i] += len(allTxcds)
				}

				// Adjust for zero-indexing
				face.VertInds[i] -= 1
				face.NormInds[i] -= 1
				face.TxcdInds[i] -= 1

				// Copy data to final arrays
				verts = append(verts,
					allVerts[face.VertInds[i]][0],
					allVerts[face.VertInds[i]][1],
					allVerts[face.VertInds[i]][2],
				)
				if face.NormInds[i] >= 0 {
					norms = append(norms,
						allNorms[face.NormInds[i]][0],
						allNorms[face.NormInds[i]][1],
						allNorms[face.NormInds[i]][2],
					)
				}
				if face.TxcdInds[i] >= 0 {
					txcds = append(txcds,
						allTxcds[face.TxcdInds[i]][0],
						allTxcds[face.TxcdInds[i]][1],
					)
				}
			}
		}

		vertCount := int32(len(group.Faces) * 3 * 3)
		model.groups = append(model.groups, modelGroup{
			DrawMode: gl.TRIANGLES,
			Start:    start,
			Count:    vertCount,
		})
		start += vertCount
	}

	gl.GenVertexArrays(1, &model.glVao)
	gl.BindVertexArray(model.glVao)
	gl.GenBuffers(3, &model.glVbos[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, model.glVbos[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(verts)*4, gl.Ptr(verts), gl.STATIC_DRAW)
	gl.VertexAttribPointer(VERT_ATTRIB, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(VERT_ATTRIB)

	if len(norms) == 0 {
		gl.DeleteBuffers(1, &model.glVbos[1])
		model.glVbos[1] = 0
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, model.glVbos[1])
		gl.BufferData(gl.ARRAY_BUFFER, len(norms)*4, gl.Ptr(norms), gl.STATIC_DRAW)
		gl.VertexAttribPointer(NORM_ATTRIB, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(NORM_ATTRIB)
	}

	if len(txcds) == 0 {
		gl.DeleteBuffers(1, &model.glVbos[2])
		model.glVbos[2] = 0
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, model.glVbos[2])
		gl.BufferData(gl.ARRAY_BUFFER, len(txcds)*4, gl.Ptr(txcds), gl.STATIC_DRAW)
		gl.VertexAttribPointer(TXCD_ATTRIB, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
		gl.EnableVertexAttribArray(TXCD_ATTRIB)
	}

	return nil
}

func (model *Model) Render(shader *Shader) {
	gl.BindVertexArray(model.glVao)

	for g := range model.groups {
		group := &model.groups[g]
		gl.DrawArrays(group.DrawMode, group.Start, group.Count)
	}
}
