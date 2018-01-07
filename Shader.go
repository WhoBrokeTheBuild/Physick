package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	glId uint32
}

func NewShader(filenames ...string) (*Shader, error) {
	glProgId := gl.CreateProgram()

	glIds := []uint32{}
	defer func() {
		for _, glId := range glIds {
			gl.DeleteShader(glId)
		}
	}()

	for _, f := range filenames {
		glId, err := compileShader(f)
		if err != nil {
			gl.DeleteProgram(glProgId)
			return nil, err
		}

		gl.AttachShader(glProgId, glId)
	}

	gl.LinkProgram(glProgId)

	var status int32
	gl.GetProgramiv(glProgId, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(glProgId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(glProgId, logLength, nil, gl.Str(log))

		return nil, fmt.Errorf("Failed to link shader program: %v", log)
	}

	shader := &Shader{
		glId: glProgId,
	}

	return shader, nil
}

func (shader *Shader) Cleanup() {
	gl.DeleteShader(shader.glId)
}

func (shader *Shader) Use() {
	gl.UseProgram(shader.glId)
}

func (shader *Shader) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(shader.glId, gl.Str(name+"\x00"))
}

func (shader *Shader) GetAttribLocation(name string) int32 {
	return gl.GetAttribLocation(shader.glId, gl.Str(name+"\x00"))
}

func compileShader(filename string) (uint32, error) {
	var shaderType uint32
	if strings.HasSuffix(filename, ".vs.glsl") {
		shaderType = gl.VERTEX_SHADER
	} else if strings.HasSuffix(filename, ".fs.glsl") {
		shaderType = gl.FRAGMENT_SHADER
	} else if strings.HasSuffix(filename, ".gs.glsl") {
		shaderType = gl.GEOMETRY_SHADER
	}

	//data, err := app.AssetFunction(filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	source := string(data) + "\x00"

	glId := gl.CreateShader(shaderType)

	glSources, free := gl.Strs(source)
	gl.ShaderSource(glId, 1, glSources, nil)
	free()

	gl.CompileShader(glId)

	var status int32
	gl.GetShaderiv(glId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(glId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(glId, logLength, nil, gl.Str(log))

		return glId, fmt.Errorf("Failed to compile shader '%v': %v", filename, log)
	}

	return glId, nil
}
