package main

import (
	"fmt"
	_ "image/png"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type Cube struct {
	vertices []float32
	size     mgl32.Vec3
	pos      mgl32.Vec3
}

func NewCube(size mgl32.Vec3, pos mgl32.Vec3) *Cube {
	var c Cube
	c.size = size
	c.pos = pos
	c.GenerateVertices()
	return &c
}

func NewRandomCube(distance float32, size float32) *Cube {
	var c Cube
	firstSize := rand.Float32() * size
	randSizeX := firstSize - rand.Float32()*firstSize/2
	randSizeY := firstSize - rand.Float32()*firstSize/2
	randSizeZ := firstSize - rand.Float32()*firstSize/2
	if randSizeX < randSizeY/2 {
		fmt.Println("firstSize", firstSize)
		fmt.Println("sizes", randSizeX, randSizeY, randSizeZ)
	}
	c.size = mgl32.Vec3{randSizeX, randSizeY, randSizeZ}
	c.pos = mgl32.Vec3{rand.Float32()*distance - distance/2, rand.Float32()*distance - distance/2, rand.Float32()*distance - distance/2}
	c.GenerateVertices()
	return &c
}

func (c *Cube) GenerateVertices() {
	// Generate a cube of size 1 around 0,0
	c.vertices = []float32{
		//  X, Y, Z, U, V
		// Bottom
		-1.0, -1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,

		// Top
		-1.0, 1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, 1.0, 0.0, 1.0,
		1.0, 1.0, 1.0, 1.0, 1.0,

		// Front
		-1.0, -1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,

		// Back
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		-1.0, 1.0, -1.0, 0.0, 1.0,
		1.0, 1.0, -1.0, 1.0, 1.0,

		// Left
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,
		-1.0, -1.0, -1.0, 0.0, 0.0,
		-1.0, -1.0, 1.0, 0.0, 1.0,
		-1.0, 1.0, 1.0, 1.0, 1.0,
		-1.0, 1.0, -1.0, 1.0, 0.0,

		// Right
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, -1.0, -1.0, 1.0, 0.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, -1.0, 1.0, 1.0, 1.0,
		1.0, 1.0, -1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0, 1.0,
	}
	// shift it by c.pos X
	for i := 0; i < len(c.vertices)/5; i++ {
		c.vertices[i*5] += c.pos[0]
		c.vertices[i*5] *= c.size[0]
	}
	// shift it by c.pos Y
	for i := 0; i < len(c.vertices)/5; i++ {
		c.vertices[i*5+1] += c.pos[1]
		c.vertices[i*5+1] *= c.size[1]
	}
	// shift it by c.pos Z
	for i := 0; i < len(c.vertices)/5; i++ {
		c.vertices[i*5+2] += c.pos[2]
		c.vertices[i*5+2] *= c.size[2]

	}
}
