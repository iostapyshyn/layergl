package main

import (
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/iostapyshyn/layergl"
	"runtime"
	"time"
)

var window *glfw.Window
var running bool = true

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

	// OpenGL version 3.3 Core
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Required for MSAA anti-aliasing
	glfw.WindowHint(glfw.Samples, 4)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)

	window, err = glfw.CreateWindow(width, height, "...", nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.SetKeyCallback(keyCallback)

	// Center window on screen
	vidmode := glfw.GetPrimaryMonitor().GetVideoMode()
	window.SetPos((vidmode.Width-width)/2, (vidmode.Height-height)/2)

	window.MakeContextCurrent()
	glfw.SwapInterval(0)

	window.Show()

	loop()
}

func loop() {
	err := layergl.Init(width, height)
	if err != nil {
		panic(err)
	}

	tex, err = layergl.NewTexture("assets/tex.png", 50, 50)
	if err != nil {
		panic(err)
	}

	bg, err := layergl.NewTexture("assets/sky.png", width, height)

	// Ticker for updating FPS
	ticker := time.NewTicker(time.Second)
	frames := 0

	defer ticker.Stop()

	worldRun()
	defer worldStop()

	for !window.ShouldClose() && running {
		layergl.Clear()

		layergl.DrawTexture(bg)

		mu.Lock()

		layergl.DrawTexture(tex.Move(570, 410))
		layergl.DrawPolygon(rect, rectColor)

		mu.Unlock()

		select {
		case <-ticker.C:
			window.SetTitle(fmt.Sprintf("%v FPS", frames))
			frames = 0
		default:
			frames++
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
