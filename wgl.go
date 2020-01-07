package wgl

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var Window *glfw.Window
var Elements []Element

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

func Loop() {
	for !Window.ShouldClose() {
		gl.ClearColor(0.2, 0.2, 0.2, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		for _, element := range Elements {
			element.Draw()
			element.Shader.CheckHotloadStatus()
		}

		Window.SwapBuffers()
		glfw.PollEvents()

	}
}
