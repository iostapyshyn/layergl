package main

import (
	"github.com/iostapyshyn/layergl"
	"math/rand"
	"sync"
	"time"
)

var mu sync.Mutex

var rect = layergl.Rectangle(layergl.Point{X: 300, Y: 200}, 200, 100)
var tex *layergl.Texture

var (
	angularVelocity = -.3
	xVelocity       = 1.2
	yVelocity       = 1.2
	rectColor       = layergl.Color{}
)

// Starts new thread calling worldUpdate() every few milliseconds
func worldRun() {
	const usec = 32768 // time in microseconds between calls of worldUpdate()

	worldInit()

	go func() {
		running = true
		var timeNow, timeLast, delta int64
		timeLast = time.Now().UnixNano() / 1000
		for running {
			timeNow = time.Now().UnixNano() / 1000
			delta += timeNow - timeLast
			timeLast = timeNow
			if delta >= usec {
				worldUpdate()
				delta -= usec
			}
		}

	}()

	for !running {
		// Wait for running flag to set true.
	}
}

// World initialization
func worldInit() {
	rand.Seed(time.Now().UnixNano())
	rectColor = layergl.Color{rand.Float64(), rand.Float64(), rand.Float64(), 0.9}
}

// Main world updating function.
func worldUpdate() {
	const texRotation = -0.2

	mu.Lock()

	rect.Move(xVelocity, yVelocity)
	rect.RotateDeg(angularVelocity)
	tex.RotateDeg(texRotation)

	if (rect.XMostPoint().X >= width && xVelocity > 0) ||
		(rect.XLeastPoint().X <= 0 && xVelocity < 0) {
		xVelocity = -xVelocity
		angularVelocity = -angularVelocity
		rectColor = layergl.Color{rand.Float64(), rand.Float64(), rand.Float64(), 0.9}
	}

	if (rect.YMostPoint().Y >= height && yVelocity > 0) ||
		(rect.YLeastPoint().Y <= 0 && yVelocity < 0) {
		yVelocity = -yVelocity
		angularVelocity = -angularVelocity
		rectColor = layergl.Color{rand.Float64(), rand.Float64(), rand.Float64(), 0.9}
	}

	mu.Unlock()
}

func worldStop() {
	running = false
}
