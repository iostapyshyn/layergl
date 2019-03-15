package main

import (
	"github.com/go-gl/glfw/v3.2/glfw"
)

// ESC to close application.
func keyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch action {
	case glfw.Release:
		switch key {
		case glfw.KeyEscape:
			w.SetShouldClose(true)
		}
	}
}
