# LayerGL
Modern OpenGL abstraction layer (In other words, graphics library) for Go.

![Minimal example window](screenshot.png?raw=true)

## Usage
### Minimal Example

```go
package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/iostapyshyn/layergl"
	"runtime"
)

var window *glfw.Window

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

	window, err = glfw.CreateWindow(width, height, "Example", nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.SetKeyCallback(keyCallback)

	// Center window on screen
	vidmode := glfw.GetPrimaryMonitor().GetVideoMode()
	window.SetPos((vidmode.Width-width)/2, (vidmode.Height-height)/2)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	window.Show()

	loop()
}

func loop() {
	renderer, err := layergl.CreateRenderer(width, height)
	if err != nil {
		panic(err)
	}

	for !window.ShouldClose() {
		renderer.Clear()

		// Draw violet triangle in the middle of a screen
		renderer.Polygon(layergl.Triangles([]layergl.Point{
			{X: width/2 - 100, Y: height/2 - 100},
			{X: width/2 + 100, Y: height/2 - 100},
			{X: width / 2, Y: height/2 + 100}}...), layergl.Color{1.0, 0.0, 1.0, 1.0})

		// Present the image
		renderer.Render()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
```

For more features, please refer to demo program source code included in the repository.

Libraries used:
 * [go-gl](https://github.com/go-gl/gl): Go bindings for OpenGL.
 * [GLFW](http://www.glfw.org): Open Source, multi-platform library for OpenGL, OpenGL ES and Vulkan development on the desktop.
 
