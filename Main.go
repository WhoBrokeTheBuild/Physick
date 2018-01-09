package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

var inputMap = map[glfw.Key]bool{
	glfw.KeyEscape: false,
	glfw.KeySpace:  false,
	glfw.KeyF2:     false,
	glfw.KeyUp:     false,
	glfw.KeyDown:   false,
	glfw.KeyLeft:   false,
	glfw.KeyRight:  false,
}

var windowSize = mgl32.Vec2{1024, 768}

var mainShader *Shader
var actors []*Actor

func init() {
	runtime.LockOSThread()
}

func main() {
	log.SetFlags(0)

	if err := glfw.Init(); err != nil {
		log.Fatalln("Failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(int(windowSize.X()), int(windowSize.Y()), "Physick", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		log.Fatalln("Failed to initialize OpenGL:", err)
	}

	log.Println("GL Version:", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Println("GL Vendor:", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Println("GLSL Version:", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	// Cornflower Blue
	gl.ClearColor(0.392156863, 0.584313725, 0.929411765, 1.0)

	mainShader, err = NewShader("assets/default.vs.glsl", "assets/default.fs.glsl")
	if err != nil {
		log.Fatalln("Failed to compile shader:", err)
	}
	mainShader.Use()

	actors = make([]*Actor, 0)

	view := mgl32.LookAtV(mgl32.Vec3{150, 150, 150}, mgl32.Vec3{0, -50, 0}, mgl32.Vec3{0, 1, 0})
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), windowSize.X()/windowSize.Y(), 1.0, 10000.0)

	gl.UniformMatrix4fv(mainShader.GetUniformLocation("u_View"), 1, false, &view[0])
	gl.UniformMatrix4fv(mainShader.GetUniformLocation("u_Projection"), 1, false, &projection[0])

	inputState := map[glfw.Key]glfw.Action{}

	frameDelay := float64(1000.0 / 60)
	frameElap := float64(0.0)
	currentFps := float32(0.0)

	fpsUpdateFrames := 0
	fpsUpdateDelay := float64(250.0)
	fpsUpdateElap := float64(0.0)

	now := func() float64 {
		return float64(time.Now().UnixNano()) / float64(time.Millisecond)
	}

	timeOffset := now()

	for !window.ShouldClose() {
		elapsedTime := now() - timeOffset
		timeOffset = now()

		for k := range inputMap {
			oldState := inputState[k]
			newState := window.GetKey(k)
			inputState[k] = newState

			inputMap[k] = oldState == glfw.Press && newState == glfw.Release
		}

		if inputMap[glfw.KeyEscape] {
			break
		}

		if inputMap[glfw.KeyF2] {
			polygonMode := int32(0)
			gl.GetIntegerv(gl.POLYGON_MODE, &polygonMode)
			if polygonMode == gl.LINE {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			}
		}

		if inputMap[glfw.KeySpace] {
			for i := 0; i < 10; i++ {
				AddCube()
			}
		}

		if inputMap[glfw.KeyLeft] {
			for i := 0; i < len(actors); i++ {
				actors[i].RigidBody.ApplyForce(mgl32.Vec3{-5.0, -5.0, 0.0}, Impulse)
			}
		}
		if inputMap[glfw.KeyRight] {
			for i := 0; i < len(actors); i++ {
				actors[i].RigidBody.ApplyForce(mgl32.Vec3{5.0, -5.0, 0.0}, Impulse)
			}
		}
		if inputMap[glfw.KeyUp] {
			for i := 0; i < len(actors); i++ {
				actors[i].RigidBody.ApplyForce(mgl32.Vec3{0.0, -5.0, -5.0}, Impulse)
			}
		}
		if inputMap[glfw.KeyDown] {
			for i := 0; i < len(actors); i++ {
				actors[i].RigidBody.ApplyForce(mgl32.Vec3{0.0, -5.0, 5.0}, Impulse)
			}
		}

		delta := float32(elapsedTime / frameDelay)

		for i := 0; i < len(actors); i++ {
			actors[i].Update(delta, float32(elapsedTime/1000.0))
			//log.Printf("%v\n", actors[i].Transform.Position)
			actors[i].Render(mainShader)

			for j := i + 1; j < len(actors); j++ {
				actors[i].RigidBody.CheckCollide(actors[j].RigidBody)
			}
		}

		frameElap += elapsedTime
		if frameDelay <= frameElap {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			for i := 0; i < len(actors); i++ {
				actors[i].Render(mainShader)
			}

			window.SwapBuffers()

			frameElap = 0.0
			fpsUpdateFrames++
		}

		fpsUpdateElap += elapsedTime
		if fpsUpdateDelay <= fpsUpdateElap {
			currentFps = float32(float64(fpsUpdateFrames)/fpsUpdateElap) * 1000.0

			title := fmt.Sprintf("Physick - %0.2f", currentFps)
			window.SetTitle(title)

			fpsUpdateElap = 0.0
			fpsUpdateFrames = 0
		}

		//time.Sleep(time.Millisecond * 16)

		glfw.PollEvents()
	}
}

func RemoveAllActors() {

}

func AddCube() {
	actor := NewActor()
	//size := (rand.Float32() * 3) + 1.0
	size := float32(1)
	model, _ := NewModelFromFile("assets/sphere.obj")
	actor.AddModel(model)
	actor.Transform.Position = mgl32.Vec3{
		rand.Float32() * 100,
		rand.Float32() * 100,
		rand.Float32() * 100,
	}
	actor.Transform.Scale = mgl32.Vec3{size, size, size}
	actor.RigidBody.Radius = size
	actor.RigidBody.Mass = size
	actor.RigidBody.ApplyForce(mgl32.Vec3{0, -9.81, 0}, Acceleration)
	actor.RigidBody.ApplyForce(mgl32.Vec3{
		(rand.Float32() - 0.5) * 10,
		(rand.Float32() - 0.5) * 10,
		(rand.Float32() - 0.5) * 10,
	}, Impulse)
	actors = append(actors, actor)
}

func DistanceSquared(p1, p2 mgl32.Vec3) float32 {
	tmp := p2.Sub(p1)
	return tmp.Dot(tmp)
}

func Distance(p1, p2 mgl32.Vec3) float32 {
	return float32(math.Sqrt(float64(DistanceSquared(p1, p2))))
}
