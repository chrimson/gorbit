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
	a := app.App()
	system := core.NewNode()

	cam := camera.New(1)
	cam.SetPosition(0, 7, 15)
	cam.LookAt(&math32.Vector3{X: 0.0, Y: 0.0, Z: 0.0}, &math32.Vector3{X: 0.0, Y: 10.0, Z: 10.0})
	camera.NewOrbitControl(cam)
	system.Add(cam)

	onResize := func(evname string, ev interface{}) {
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	sunShape := geometry.NewSphere(2, 360, 360)
	sunColor := material.NewStandard(&math32.Color{R: 1.0, G: 0.9, B: 0.7})
	sunObj := graphic.NewMesh(sunShape, sunColor)
	sun := core.NewNode()
	sun.Add(sunObj)
	sun.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.8))
	system.Add(sun)

	earthShape := geometry.NewSphere(0.5, 360, 360)
	earthImage := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	earthTexture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	earthTexture.SetShininess(10)
	earthTexture.AddTexture(earthImage("earth_clouds_big.jpg"))
	earthTilt := graphic.NewMesh(earthShape, earthTexture)
	earthTilt.RotateZ(23.4 * math32.Pi / 180)

	earthAxisGeometry := geometry.NewGeometry()
	earthAxisVertices := math32.NewArrayF32(0, 0)
	earthAxisVertices.Append(
		0.0, 1.0, 0.0,
		0.0, -1.0, 0.0,
	)
	earthAxisGeometry.AddVBO(gls.NewVBO(earthAxisVertices).AddAttrib(gls.VertexPosition))
	earthAxisMaterial := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	earthAxis := graphic.NewLines(earthAxisGeometry, earthAxisMaterial)
	earthTilt.Add(earthAxis)

	moonShape := geometry.NewSphere(0.15, 360, 360)
	moonImage := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	moonTexture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	moonTexture.SetShininess(10)
	moonTexture.AddTexture(moonImage("moon.jpg"))

	moonDistance := graphic.NewMesh(moonShape, moonTexture)
	moonDistance.TranslateX(1.5)
	moon := core.NewNode()
	moon.Add(moonDistance)

	earthDistance := core.NewNode()
	earthDistance.Add(moon)
	earthDistance.Add(earthTilt)
	earthDistance.TranslateX(10.0)
	earth := core.NewNode()
	earth.Add(earthDistance)
	system.Add(earth)

	earthPathCircle := geometry.NewGeometry()
	earthPathPoints := math32.NewArrayF32(0, 0)
	for x := float32(-1.0); x < 1.0; x = x + 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		earthPathPoints.Append(10.0*x, 0.0, 10.0*z)
	}
	for x := float32(1.0); x > -1.0; x = x - 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		earthPathPoints.Append(10.0*x, 0.0, -10.0*z)
	}
	earthPathCircle.AddVBO(gls.NewVBO(earthPathPoints).AddAttrib(gls.VertexPosition))
	earthPathMaterial := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	earthPath := graphic.NewLineStrip(earthPathCircle, earthPathMaterial)
	system.Add(earthPath)

	lights := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 20.0)
	lights.SetPosition(0, 0, 0)
	system.Add(lights)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		delta := float32(deltaTime.Seconds()) * 2 * math32.Pi / 365
		earth.RotateY(delta)
		earthDistance.RotateY(-delta)
		earthTilt.RotateY(delta * 365)
		moon.RotateY(delta * 365 / 27)

		renderer.Render(system, cam)
	})
}
