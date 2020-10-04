package main

import (
	"fmt"
	_ "image/png"
	"math"

	// "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/thabaptiser/goclid/gl"
)

type Action int

const (
	PLAYER_FORWARD   Action = iota
	PLAYER_BACKWARD  Action = iota
	PLAYER_LEFT      Action = iota
	PLAYER_RIGHT     Action = iota
	PLAYER_LOOKUP    Action = iota
	PLAYER_LOOKDOWN  Action = iota
	PLAYER_LOOKLEFT  Action = iota
	PLAYER_LOOKRIGHT Action = iota
	PROGRAM_QUIT     Action = iota
)

type InputManager struct {
	actionToKeyMap map[Action]glfw.Key
	keysPressed    [glfw.KeyLast]bool

	firstCursorAction    bool
	cursor               mgl32.Vec2
	cursorChange         mgl32.Vec2
	cursorLast           mgl32.Vec2
	bufferedCursorChange mgl32.Vec2
}

func NewInputManager() *InputManager {
	actionToKeyMap := map[Action]glfw.Key{
		PLAYER_FORWARD:   glfw.KeyW,
		PLAYER_BACKWARD:  glfw.KeyS,
		PLAYER_LEFT:      glfw.KeyA,
		PLAYER_RIGHT:     glfw.KeyD,
		PLAYER_LOOKUP:    glfw.KeyI,
		PLAYER_LOOKDOWN:  glfw.KeyK,
		PLAYER_LOOKLEFT:  glfw.KeyJ,
		PLAYER_LOOKRIGHT: glfw.KeyL,
		PROGRAM_QUIT:     glfw.KeyEscape,
	}

	return &InputManager{
		actionToKeyMap:    actionToKeyMap,
		firstCursorAction: false,
	}
}

// IsActive returns whether the given Action is currently active
func (im *InputManager) IsActive(a Action) bool {
	return im.keysPressed[im.actionToKeyMap[a]]
}

func (im *InputManager) keyCallback(window *glfw.Window, key glfw.Key, scancode int,
	action glfw.Action, mods glfw.ModifierKey) {

	// timing for key events occurs differently from what the program loop requires
	// so just track what key actions occur and then access them in the program loop
	switch action {
	case glfw.Press:
		im.keysPressed[key] = true
	case glfw.Release:
		im.keysPressed[key] = false
	}
}
func inputActions(im *InputManager, program uint32, cc *CameraCoords, elapsed float32) {
	// fmt.Println("In inputActions")
	// Move the camera position
	moveSpeed := elapsed * 1000
	cameraSpeed := elapsed
	if im.IsActive(PLAYER_FORWARD) {
		cc.pos = cc.pos.Add(cc.front.Mul(moveSpeed))
	}
	if im.IsActive(PLAYER_BACKWARD) {
		cc.pos = cc.pos.Sub(cc.front.Mul(moveSpeed))
	}
	if im.IsActive(PLAYER_RIGHT) {
		cc.pos = cc.pos.Add(cc.front.Cross(UpVector()).Normalize().Mul(moveSpeed))
	}
	if im.IsActive(PLAYER_LEFT) {
		cc.pos = cc.pos.Sub(cc.front.Cross(UpVector()).Normalize().Mul(moveSpeed))
	}
	// wrap the movement around
	cc.filterPosition()

	// move the camera angle
	if im.IsActive(PLAYER_LOOKRIGHT) {
		cc.yaw += cameraSpeed / 2
	}
	if im.IsActive(PLAYER_LOOKLEFT) {
		cc.yaw -= cameraSpeed / 2
	}
	if im.IsActive(PLAYER_LOOKUP) {
		cc.pitch += cameraSpeed / 2
	}
	if im.IsActive(PLAYER_LOOKDOWN) {
		cc.pitch -= cameraSpeed / 2
	}
	if cc.yaw > math.Pi*2 {
		cc.yaw -= math.Pi * 2
	}
	if cc.yaw < 0 {
		cc.yaw += math.Pi * 2
	}
	if cc.pitch > math.Pi*2 {
		cc.pitch -= math.Pi * 2
	}
	if cc.pitch < 0 {
		cc.pitch += math.Pi * 2
	}

	//If the pitch is > 90 degrees, but less then 180, we need to flip the yaw.
	yaw := cc.yaw
	if cc.pitch > math.Pi*1.5 {
		yaw = cc.yaw + math.Pi
	}
	if cc.pitch < math.Pi/2 {
		yaw = cc.yaw + math.Pi
	}

	// // fmt.Println("pitch", cc.pitch)
	// fmt.Println("yaw", yaw)

	// IDEA: change UpVector while moving the camera, that way
	// turning "right" always works
	cc.front = mgl32.Vec3{
		float32(math.Cos(float64(yaw)) * math.Cos(float64(cc.pitch))),
		float32(math.Sin(float64(cc.pitch))),
		float32(math.Sin(float64(yaw)) * math.Cos(float64(cc.pitch))),
	}.Normalize()
	// fmt.Println("cc.front", cc.front)

	camera := mgl32.LookAtV(cc.pos, cc.pos.Add(cc.front), UpVector())
	fmt.Println("Camera", camera)
	// fmt.Println("Pos", cc.pos)
	// fmt.Println("Front", cc.front)
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])
}
