package main

import (
	"fmt"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/texture"
	"github.com/g3n/engine/window"

	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
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

	s1 := gui.NewHSlider(400, 32)
	s1.SetPosition(10, 10)
	s1.SetValue(0.2)
	s1.SetText(fmt.Sprintf("%1.2f", s1.Value()))
	s1.Subscribe(gui.OnChange, func(evname string, ev interface{}) {
		// a.Log().Info("Slider 1 OnChange: %v", s1.Value())
		s1.SetText(fmt.Sprintf("%1.2f", s1.Value()))
	})
	system.Add(s1)

	sunShape := geometry.NewSphere(2, 360, 360)
	sunColor := material.NewStandard(&math32.Color{R: 1.0, G: 0.9, B: 0.7})
	sunObj := graphic.NewMesh(sunShape, sunColor)
	sun := core.NewNode()
	sun.Add(sunObj)
	sun.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.3))
	system.Add(sun)

	earthShape := geometry.NewSphere(0.5, 360, 360)
	earthImage := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	earthTexture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	earthTexture.SetShininess(10)
	earthTexture.AddTexture(earthImage("earth.jpg"))
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

	moonPathCircle := geometry.NewGeometry()
	moonPathPoints := math32.NewArrayF32(0, 0)
	for x := float32(-1.0); x < 1.0; x = x + 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		moonPathPoints.Append(1.5*x, 0.0, 1.5*z)
	}
	for x := float32(1.0); x > -1.0; x = x - 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		moonPathPoints.Append(1.5*x, 0.0, -1.5*z)
	}
	moonPathCircle.AddVBO(gls.NewVBO(moonPathPoints).AddAttrib(gls.VertexPosition))
	moonPathMaterial := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	moonPath := graphic.NewLineStrip(moonPathCircle, moonPathMaterial)

	moonPlane := core.NewNode()
	moonPlane.Add(moonDistance)
	moonPlane.Add(moonPath)
	moonPlane.RotateZ(5.14 * math32.Pi / 180)
	moon := core.NewNode()
	moon.Add(moonPlane)

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

	lights := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 30.0)
	lights.SetPosition(4.0, 3.0, -2.0)
	system.Add(lights)

	a.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		a.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		delta := float32(deltaTime.Seconds()) * 2 * math32.Pi / 365.0
		earth.RotateY(delta)
		earthDistance.RotateY(-delta)
		earthTilt.RotateY(delta * 365.0)
		moonPlane.RotateY(delta * 365.0 / 27.3)

		renderer.Render(system, cam)
	})
}
