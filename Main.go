package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var WindowSize = mgl32.Vec2{1024, 768}

func init() {
	runtime.LockOSThread()
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(int(WindowSize.X()), int(WindowSize.Y()), "Physick", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println("GL Version:", gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println("GL Vendor:", gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println("GLSL Version:", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

    // Cornflower Blue
	gl.ClearColor(0.392156863, 0.584313725, 0.929411765, 1.0)

	shader, err := NewShader("assets/default.vs.glsl", "assets/default.fs.glsl")
	if err != nil {
		panic(err)
	}
	shader.Use()

	actor := NewActor()
	actor.AddModel(NewCubeModel(shader, 1.0))

	view := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), WindowSize.X()/WindowSize.Y(), 0.1, 100.0)

	gl.UniformMatrix4fv(shader.GetUniformLocation("u_View"), 1, false, &view[0])
	gl.UniformMatrix4fv(shader.GetUniformLocation("u_Projection"), 1, false, &projection[0])

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		actor.Render(shader)

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
