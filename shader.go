package wogl

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v4.1-core/gl"
)

type ShaderProgram struct {
	programID   uint32
	vertID      uint32
	fragID      uint32
	vertPath    string
	fragPath    string
	attributes  map[string]int32
	uniforms    map[string]int32
	vertModTime time.Time
	fragModTime time.Time
}

func NewShader(vertexSourcePath, fragmentSourcePath string) *ShaderProgram {
	var shaderProgram ShaderProgram

	shaderProgram.attributes = make(map[string]int32)
	shaderProgram.uniforms = make(map[string]int32)
	shaderProgram.vertPath = vertexSourcePath
	shaderProgram.fragPath = fragmentSourcePath

	vertexShaderSource := loadShader(vertexSourcePath)
	fragmentShaderSource := loadShader(fragmentSourcePath)

	shaderProgram.vertModTime = checkModifiedTime(shaderProgram.vertPath)
	shaderProgram.fragModTime = checkModifiedTime(shaderProgram.fragPath)

	shaderProgram.vertID = compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	shaderProgram.fragID = compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)

	shaderProgram.programID = createProgram(shaderProgram.vertID, shaderProgram.fragID)
	shaderProgram.LoadAttributes()
	shaderProgram.LoadUniforms()
	return &shaderProgram
}

func (shaderProgram *ShaderProgram) Use() {
	gl.UseProgram(shaderProgram.programID)
}

func (shader *ShaderProgram) SetFloat(name string, f float32) {
	nameCstr := gl.Str(name + "\x00")
	location := gl.GetUniformLocation(shader.programID, nameCstr)
	gl.Uniform1f(location, f)
}

func (shader *ShaderProgram) SetVec2(name string, v mgl32.Vec2) {
	location := gl.GetUniformLocation(shader.programID, gl.Str(name+"\x00"))
	gl.Uniform2f(location, v.X(), v.Y())
}

func (shader *ShaderProgram) SetMatrix4x4(name string, m mgl32.Mat4) {
	gl.UniformMatrix4fv(shader.GetUniformLocation(name), 1, false, &m[0])
}

func (shader *ShaderProgram) CheckHotloadStatus() {
	vertNewTime := checkModifiedTime(shader.vertPath)
	fragNewTime := checkModifiedTime(shader.fragPath)

	if shader.vertModTime.Before(vertNewTime) || shader.fragModTime.Before(fragNewTime) {
		fmt.Println("Reloading vertex and fragment shader")
		shader.vertModTime = vertNewTime
		shader.vertID = compileShader(loadShader(shader.vertPath), gl.VERTEX_SHADER)
		shader.fragModTime = fragNewTime
		shader.fragID = compileShader(loadShader(shader.fragPath), gl.FRAGMENT_SHADER)

		gl.DeleteProgram(shader.programID)
		shader.programID = createProgram(shader.vertID, shader.fragID)
	}
}

func createProgram(vertID, fragID uint32) uint32 {
	programID := gl.CreateProgram()
	gl.AttachShader(programID, vertID)
	gl.AttachShader(programID, fragID)
	gl.LinkProgram(programID)

	gl.DeleteShader(vertID)
	gl.DeleteShader(fragID)
	return programID
}

func checkModifiedTime(filePath string) time.Time {
	file, err := os.Stat(filePath)
	if err != nil {
		fmt.Println(err)
	}
	return file.ModTime()
}

func loadShader(source string) string {
	b, err := ioutil.ReadFile(source)
	if err != nil {
		fmt.Printf("%v", err)
	}
	s := string(b)
	return s + "\x00"
}

func compileShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		fmt.Printf("compileShader: failed to compile %v: %v", source, log)
		return 0
	}

	return shader
}

func (program *ShaderProgram) LoadUniforms() {
	var count int32
	gl.GetProgramiv(program.programID, gl.ACTIVE_UNIFORMS, &count)

	var i int32
	for i = 0; i < count; i++ {
		var length int32
		var size int32
		var xtype uint32
		var name [30]uint8
		gl.GetActiveUniform(program.programID, uint32(i), 30, &length, &size, &xtype, &name[0])

		program.uniforms[gl.GoStr(&name[0])] = gl.GetUniformLocation(program.programID, &name[0])
	}
	fmt.Printf("Active uniforms: %d\n", count)
	for name, i := range program.uniforms {
		fmt.Printf("%5d: %s\n", i, name)
	}
}

func (program *ShaderProgram) LoadAttributes() {
	var count int32
	gl.GetProgramiv(program.programID, gl.ACTIVE_ATTRIBUTES, &count)

	var i int32
	for i = 0; i < count; i++ {
		var length int32
		var size int32
		var xtype uint32
		var name [30]uint8
		gl.GetActiveAttrib(program.programID, uint32(i), 30, &length, &size, &xtype, &name[0])

		program.attributes[gl.GoStr(&name[0])] = gl.GetAttribLocation(program.programID, &name[0])
	}
	fmt.Printf("Active attributes: %d\n", count)
	for name, i := range program.attributes {
		fmt.Printf("%5d: %s\n", i, name)
	}
}

func (program *ShaderProgram) GetUniformLocation(name string) int32 {
	return program.uniforms[name]
}
