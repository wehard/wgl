package wgl

import (
	"fmt"
	"runtime"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var Window *glfw.Window
var Elements []Element
var ClearColor mgl32.Vec4

func Init(width, height int, title string) func() {
	err, cleanupFunc := initGlfw(width, height, title)
	if err != nil {
		panic(err)
	}
	err = initGl()
	if err != nil {
		panic(err)
	}
	return cleanupFunc
}

func initGlfw(width, height int, title string) (error, func()) {

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
	glfw.WindowHint(glfw.CenterCursor, 1)
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Created window")
	Window = window

	Window.MakeContextCurrent()

	cleanup := func() {
		glfw.Terminate()
	}
	return nil, cleanup
}

func initGl() error {
	err := gl.Init()
	if err != nil {
		return err
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)
	return nil
}

func AddElement(element *Element) {
	Elements = append(Elements, *element)
}

func SetKeyCallback(cbfunc func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey)) {
	Window.SetKeyCallback(cbfunc)
}

func Loop(onUpdate, onRender func()) {
	for !Window.ShouldClose() {
		gl.ClearColor(ClearColor[0], ClearColor[1], ClearColor[2], ClearColor[3])
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		onUpdate()
		for _, element := range Elements {
			element.Draw()
			element.Shader.CheckHotloadStatus()
		}
		onRender()

		Window.SwapBuffers()
		glfw.PollEvents()

	}
}
