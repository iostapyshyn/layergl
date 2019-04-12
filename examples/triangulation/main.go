package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/iostapyshyn/layergl"
	"log"
	"runtime"
)

var window *glfw.Window

var polygon = new(layergl.VertexObject)
var wireframe bool = true
var mouseHeld bool = false

const (
	width  = 640
	height = 480
)

var (
	bgColor      layergl.Color = layergl.Color{0.95, 0.95, 0.95, 1.0}
	polygonColor layergl.Color = layergl.Color{1.0, 0.0, 0.0, 0.5}
	wireColor    layergl.Color = layergl.Color{0.2, 0.2, 0.2, 1.0}
)

func init() {
	runtime.LockOSThread()
}

func keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Release:
		switch key {
		case glfw.KeyEscape:
			window.SetShouldClose(true)
		}

	case glfw.Press:
		switch key {
		case glfw.KeySpace: // Triangulate the polygon.
			if len(polygon.Indices) == 0 {
				if err := polygon.Triangulate(); err != nil {
					log.Println(err)
				} else {
					log.Printf("%v vertices triangulated into %v triangles.", len(polygon.Vertices), len(polygon.Indices)/3)
				}
			} else {
				polygon.Indices = polygon.Indices[:0]
			}
		case glfw.KeyW: // Toggle wireframe.
			wireframe = !wireframe
		case glfw.KeyC: // Clear.
			polygon.Vertices = polygon.Vertices[:0]
			polygon.Indices = polygon.Indices[:0]
		case glfw.KeyZ: // Undo.
			if len(polygon.Vertices) > 0 {
				polygon.Vertices = append(polygon.Vertices[:len(polygon.Vertices)-1], polygon.Vertices[len(polygon.Vertices):]...)
				polygon.Indices = polygon.Indices[:0]
			}
		}

	case glfw.Repeat:
		switch key {
		case glfw.KeyZ: // Undo.
			polygon.Vertices = append(polygon.Vertices[:len(polygon.Vertices)-1], polygon.Vertices[len(polygon.Vertices):]...)
			polygon.Indices = polygon.Indices[:0]
		}
	}
}

func mouseCallback(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if button == glfw.MouseButtonLeft && action == glfw.Press {
		x, y := window.GetCursorPos()
		polygon.Vertices = append(polygon.Vertices, layergl.Point{X: x, Y: height - y})
		polygon.Indices = polygon.Indices[:0]

		mouseHeld = true
	} else if button == glfw.MouseButtonLeft && action == glfw.Release {
		mouseHeld = false
	}
}

func main() {
	var err error
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	// OpenGL version 3.3 Core.
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Required for MSAA anti-aliasing.
	glfw.WindowHint(glfw.Samples, 4)

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Visible, glfw.False)

	window, err = glfw.CreateWindow(width, height, "Polygon Triangulation", nil, nil)
	if err != nil {
		panic(err)
	}

	defer window.Destroy()

	window.SetKeyCallback(keyCallback)
	window.SetMouseButtonCallback(mouseCallback)

	// Center window on the screen.
	vidmode := glfw.GetPrimaryMonitor().GetVideoMode()
	window.SetPos((vidmode.Width-width)/2, (vidmode.Height-height)/2)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	window.Show()

	loop()
}

func loop() {
	err := layergl.Init(width, height)
	if err != nil {
		panic(err)
	}

	layergl.ClearColor(bgColor)

	for !window.ShouldClose() {
		// Add points if mouse if being dragged.
		if mouseHeld {
			mouseDrag()
		}

		layergl.Clear()

		// Draw triangulated polygon or preview line if polygon is not in triangulated state.
		if len(polygon.Indices) >= 3 {
			layergl.DrawVertexObject(polygon, polygonColor)
		} else if len(polygon.Vertices) > 0 {
			for _, p := range polygon.Vertices {
				layergl.DrawPoint(p, 2, wireColor)
			}

			layergl.DrawLines(append(polygon.Vertices, polygon.Vertices[0]), wireColor)
		}

		if wireframe {
			drawWireframe()
		}

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func mouseDrag() {
	x, y := window.GetCursorPos()

	// Add new point only if the distance from the previous one is greater than 5px
	if len(polygon.Vertices) != 0 &&
		(layergl.Distance(polygon.Vertices[len(polygon.Vertices)-1], layergl.Point{x, height - y}) > 5) {

		polygon.Vertices = append(polygon.Vertices, layergl.Point{X: x, Y: height - y})
		polygon.Indices = polygon.Indices[:0]
	}
}

func drawWireframe() {
	for i := 0; i < len(polygon.Indices); i += 3 {
		layergl.DrawLines([]layergl.Point{
			polygon.Vertices[polygon.Indices[i]],
			polygon.Vertices[polygon.Indices[i+1]],
			polygon.Vertices[polygon.Indices[i+2]],
			polygon.Vertices[polygon.Indices[i]],
		}, wireColor)
	}
}
