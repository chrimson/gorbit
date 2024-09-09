package main

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
)

func createMoon() *core.Node {
	moonShape := geometry.NewSphere(0.15, 360, 360)
	moonImage := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	moonTexture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	moonTexture.SetShininess(5)
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
	moon := core.NewNode()
	moon.Add(moonDistance)
	moon.Add(moonPath)
	moon.RotateZ(5.14 * math32.Pi / 180)

	return moon
}
