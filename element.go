package wgl

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	x, y, z    float64
	r, g, b, a float64
	u, v       float64
}

// Element holds GL ids, vertex data, transform and shader
type Element struct {
	vaoID     uint32
	vboID     uint32
	eboID     uint32
	vertices  []float32
	indices   []uint32
	uvs       []float32
	Transform Transform
	data      []Vertex
	Shader    *ShaderProgram
}

// NewElement creates a new instance of element, generates it's buffers and returns a pointer to it
func NewElement(shader *ShaderProgram) *Element {
	//var element Element
	//element.vertices = []float32{
	//	0, 0.5, 0,
	//	-0.5, -0.5, 0,
	//	0.5, -0.5, 0,
	//}
	element := MakeQuad()
	element.genBuffers()
	element.Transform = NewTransform()
	element.Shader = shader
	return element
}

func MakeQuad() *Element {
	var quad Element
	quad.vertices = []float32{
		-0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0, 0.0,
		-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0,
		0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 1.0, 1.0, 1.0,
		0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 1.0, 1.0, 0.0,
	}
	quad.indices = []uint32{
		0, 1, 2,
		2, 3, 0,
	}
	quad.uvs = []float32{
		0.0, 0.0,
		0.0, 1.0,
		1.0, 1.0,
		1.0, 0.0,
	}
	quad.data = []Vertex{
		{x: -0.5, y: 0.5, z: 0.0, r: 1.0, g: 0.0, b: 0.0, a: 1.0},
		{x: -0.5, y: -0.5, z: 0.0, r: 0.0, g: 1.0, b: 0.0, a: 1.0},
		{x: 0.5, y: -0.5, z: 0.0, r: 0.0, g: 0.0, b: 1.0, a: 1.0},
		{x: 0.5, y: 0.5, z: 0.0, r: 1.0, g: 0.0, b: 1.0, a: 1.0},
	}
	return &quad
}

func (element *Element) genBuffers() {

	gl.GenBuffers(1, &element.eboID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, element.eboID)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(element.indices), gl.Ptr(element.indices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &element.vboID)
	gl.BindBuffer(gl.ARRAY_BUFFER, element.vboID)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(element.vertices), gl.Ptr(element.vertices), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &element.vaoID)
	gl.BindVertexArray(element.vaoID)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 9*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, 9*4, gl.PtrOffset(5*4))
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)
}

// Draw element using it's shader
func (element *Element) Draw() {

	element.Shader.Use()

	rotateX := mgl32.Rotate3DX(mgl32.DegToRad(element.Transform.Rotation.X()))
	rotateY := mgl32.Rotate3DY(mgl32.DegToRad(element.Transform.Rotation.Y()))
	rotateZ := mgl32.Rotate3DZ(mgl32.DegToRad(element.Transform.Rotation.Z()))

	pos := element.Transform.Position
	modelMatrix := mgl32.Translate3D(pos[0], pos[1], pos[2]).Mul4(rotateX.Mul3(rotateY).Mul3(rotateZ).Mat4())
	element.Shader.SetMatrix4x4("model_matrix", modelMatrix)

	gl.BindVertexArray(element.vaoID)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, element.eboID)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
}
