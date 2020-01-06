package wogl

import (
	"runtime"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.0/glfw"
)

type WoglData struct {
	Window glfw.Window
	Elements []Element
}

func NewWogl() *WoglData {
	var data WoglData

	return &data
}

func (data *WoglData) InitGlfw() (error, func()) {

	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	data.Window, err := glfw.CreateWindow(1600, 1600, "go_fractals", nil, nil)
	if err != nil {
		panic(err)
	}

	data.Window.MakeContextCurrent()

	cleanup := func() {
		glfw.Terminate()
	}
	return nil, cleanup
}

func (data *WoglData) InitGl() {
	err = gl.Init()
	if err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
}

func (data *WoglData) AddElement(element *Element) {
	data.Elements = append(data.Elements, element)
}

func (data *WoglData) Loop() {
	for !data.Window.ShouldClose() {
		gl.ClearColor(0.2, 0.2, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, e := range data.Elements {
			element.Draw()
		}

		window.SwapBuffers()
		glfw.PollEvents()
		shader.CheckHotloadStatus()
	}
}