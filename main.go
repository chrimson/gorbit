package main

import (
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/texture"

	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"

	"github.com/g3n/engine/window"
)

func main() {
	// Create application and scene
	a := app.App()
	systemOrbits := core.NewNode()
	earthOrbit := core.NewNode()

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 10)
	systemOrbits.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	earthShape := geometry.NewSphere(0.5, 360, 360)

	earthImage := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	earthTexture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	earthTexture.SetShininess(1)
	earthTexture.AddTexture(earthImage("earth_clouds_big.jpg"))

	earth := graphic.NewMesh(earthShape, earthTexture)
	earth.TranslateX(5.0)
	earth.RotateZ(0.4084)
	earthOrbit.Add(earth)
	systemOrbits.Add(earthOrbit)

	systemOrbits.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.2))
	sunLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 20.0)
	sunLight.SetPosition(0, 0, 0)
	systemOrbits.Add(sunLight)

	// Run the application
	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		delta := float32(deltaTime.Seconds()) * 2 * math32.Pi / 20
		earthOrbit.RotateY(delta)
		earth.RotateY(delta * 10)

		renderer.Render(systemOrbits, cam)
	})
}
