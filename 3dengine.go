package main

import (
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"strings"

	// "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thabaptiser/goclid/gl"
)

// import "fmt"
//
// func main() {
// 	fmt.Println("vim-go")
// }

// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.

// Action is a configurable abstraction of a key press

func UpVector() mgl32.Vec3 {
	return mgl32.Vec3{0, 1, 0}
}

const windowWidth = 1000
const windowHeight = 800
const gameSize = 100

func init() {
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}

type CameraCoords struct {
	pos   mgl32.Vec3
	front mgl32.Vec3
	yaw   float32
	pitch float32
	// x               float32
	// y               float32
	// z               float32
	// horizontalAngle float32
	// verticalAngle   float32
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	im := NewInputManager()
	var cc CameraCoords
	cc.pos = mgl32.Vec3{6, 6, 6}
	cc.front = mgl32.Vec3{1, 0, 0}
	cc.yaw = -2
	cc.pitch = -0.5

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	window.SetKeyCallback(im.keyCallback)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 100000.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(cc.pos, cc.front, UpVector())
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	texture, err := newTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Generate the cubes
	var verticies []float32
	distance := float32(gameSize * 5)
	size := float32(0.3)
	for i := 0; i < 400; i++ {
		if i%1000 == 0 {
			fmt.Println("cubes generated: ", i)
		}
		cube := NewRandomCube(distance, size)
		verticies = append(verticies, cube.vertices...)
	}
	finalVertices := verticesToCopiedVerticiesInEveryDirection(verticies, 4)
	// finalVertices := verticies
	// cube2 := NewCube(1, mgl32.Vec3{0, 0, 0})
	// finalVertices := append(append([]float32{}, cube1.vertices...), cube2.vertices...)
	fmt.Println("Printing cubes, vertices array length: ", len(finalVertices))

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(finalVertices)*4, gl.Ptr(finalVertices), gl.STATIC_DRAW)
	// gl.BufferData(gl.ARRAY_BUFFER, len(triangleVertices)*3, gl.Ptr(triangleVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	// gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.ClearColor(0, 0, 0, 0)
	// Try out Fog
	// gl.Fogi(GL_FOG_MODE, gl.GL_LINEAR) // Fog Mode
	// gl.Fogfv(GL_FOG_COLOR, fogColor)   // Set Fog Color
	// gl.Fogf(GL_FOG_DENSITY, 0.35)      // How Dense Will The Fog Be
	// gl.Hint(GL_FOG_HINT, GL_DONT_CARE) // Fog Hint Value
	// gl.Fogf(GL_FOG_START, 1.0)         // Fog Start Depth
	// gl.Fogf(GL_FOG_END, 5.0)           // Fog End Depth
	// gl.Enable(GL_FOG)                  // Enables GL_FOG

	angle := 0.0
	previousTime := glfw.GetTime()
	fmt.Println("OpenGL version", version)

	for !window.ShouldClose() {
		// fmt.Println("In main loop")
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed
		// model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		// Handle Input
		inputActions(im, program, &cc, float32(elapsed))

		// Render
		gl.UseProgram(program)
		// gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)

		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(finalVertices)/5))

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func (cc *CameraCoords) filterPosition() {
	// wrap the movement around
	for i := 0; i < 3; i++ {
		if cc.pos[i] > gameSize {
			cc.pos[i] -= 2 * float32(gameSize)
		}
		if cc.pos[i] < -gameSize {
			cc.pos[i] += 2 * float32(gameSize)
		}
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
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

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

//IDEA: use this shader to return a gl_Position that wraps?
var vertexShader = `
#version 330
#define FOG_START 100
#define FOG_END 1000
uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;
in vec3 vert;
in vec2 vertTexCoord;
out vec2 fragTexCoord;
out float fogAmount;

float fogFactorLinear(
  const float dist,
  const float start,
  const float end
) {
  return 1.0 - clamp((end - dist) / (end - start), 0.0, 1.0);
}

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
      float fogDistance = length(gl_Position.xyz);
      fogAmount = fogFactorLinear(fogDistance, FOG_START, FOG_END);
}
` + "\x00"

var fragmentShader = `
#version 330
uniform sampler2D tex;
uniform vec4 u_fogColor;
in vec2 fragTexCoord;
in float fogAmount;
out vec4 outputColor;
void main() {
    outputColor = texture(tex, fragTexCoord);
  outputColor = mix(outputColor, u_fogColor, fogAmount);  
}
` + "\x00"

func verticesToCopiedVerticiesInEveryDirection(vertices []float32, depth int) []float32 {
	var retVertices []float32
	for x := -depth; x < depth+1; x++ {
		for y := -depth; y < depth+1; y++ {
			for z := -depth; z < depth+1; z++ {
				for i := 0; i < len(vertices)/5; i++ {
					newVertex := vertices[i*5 : i*5+5]
					retVertices = append(retVertices,
						newVertex[0]+float32(x*gameSize),
						newVertex[1]+float32(y*gameSize),
						newVertex[2]+float32(z*gameSize),
						newVertex[3],
						newVertex[4])
				}
			}
		}
	}
	// fmt.Println("before", vertices)
	// fmt.Println("after", retVertices)
	return retVertices
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/thabaptiser/goclid")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}
