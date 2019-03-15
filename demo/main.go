package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/iostapyshyn/glass"
	"math/rand"
	"runtime"
	"time"
)

var window *glfw.Window
var running bool = true

var rect = glass.Rectangle(glass.Point{X: 200, Y: 300}, 200, 100)
var tex = glass.Texture{}

const (
	width  = 640
	height = 480
)

func init() {
	runtime.LockOSThread()
}

func main() {
	var err error
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)

	// Required for MSAA anti-aliasing
	glfw.WindowHint(glfw.Samples, 4)

	window, err = glfw.CreateWindow(width, height, "Title", nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.SetKeyCallback(keyCallback)

	vidmode := glfw.GetPrimaryMonitor().GetVideoMode()
	window.SetPos((vidmode.Width-width)/2, (vidmode.Height-height)/2)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	window.Show()

	loop()
}

func physicsRun() {
	const usec = 32768 // time in microseconds between calls of physicsUpdate()
	rand.Seed(time.Now().UnixNano())
	rectColor = glass.Color{rand.Float32(), rand.Float32(), rand.Float32(), 1.0}

	go func() {
		running = true
		var timeNow, timeLast, delta int64
		timeLast = time.Now().UnixNano() / 1000
		for running {
			timeNow = time.Now().UnixNano() / 1000
			delta += timeNow - timeLast
			timeLast = timeNow
			if delta >= usec {
				physicsUpdate()
				delta -= usec
			}
		}

	}()

	for !running {
		// Wait for running flag to set true
	}
}

var (
	angularVelocity = float64(0.3)
	xVelocity       = float32(1.0)
	yVelocity       = float32(1.0)
	rectColor       = glass.Color{}
)

func physicsUpdate() {
	rect = rect.Move(xVelocity, yVelocity).RotateDeg(angularVelocity)

	if (rect.XMostPoint().X >= width && xVelocity > 0) ||
		(rect.XLeastPoint().X <= 0 && xVelocity < 0) {
		xVelocity = -xVelocity
		angularVelocity = -angularVelocity
		rectColor = glass.Color{rand.Float32(), rand.Float32(), rand.Float32(), 1.0}
	}

	if (rect.YMostPoint().Y >= height && yVelocity > 0) ||
		(rect.YLeastPoint().Y <= 0 && yVelocity < 0) {
		yVelocity = -yVelocity
		angularVelocity = -angularVelocity
		rectColor = glass.Color{rand.Float32(), rand.Float32(), rand.Float32(), 1.0}
	}
}

func physicsStop() {
	running = false
}

func loop() {
	renderer, err := glass.CreateRenderer(width, height)
	if err != nil {
		panic(err)
	}

	renderer.ClearColor(glass.Color{0.15, 0.1, 0.2, 1.0})

	tex, err = glass.NewTexture("tex.png", 50, 50)
	if err != nil {
		panic(err)
	}

	bg, err := glass.NewTexture("sky.png", 640, 480)

	physicsRun()

	for !window.ShouldClose() && running {
		renderer.Clear()

		renderer.Texture(bg)

		renderer.Polygon(rect, rectColor)
		renderer.Texture(tex.Move(570, 410))

		renderer.Render()

		tex = tex.RotateDeg(-0.2)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	physicsStop()
}
