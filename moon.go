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

func newMoon() *core.Node {
	shape := geometry.NewSphere(0.15, 360, 360)
	image := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	texture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	texture.SetShininess(5)
	texture.AddTexture(image("moon.jpg"))
	side := graphic.NewMesh(shape, texture)
	side.SetRotationY(math32.Pi)
	distance := core.NewNode()
	distance.Add(side)
	distance.TranslateX(1.5)

	pathCircle := geometry.NewGeometry()
	pathPoints := math32.NewArrayF32(0, 0)
	for x := float32(-1.0); x < 1.0; x = x + 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		pathPoints.Append(1.5*x, 0.0, 1.5*z)
	}
	for x := float32(1.0); x > -1.0; x = x - 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		pathPoints.Append(1.5*x, 0.0, -1.5*z)
	}
	pathCircle.AddVBO(gls.NewVBO(pathPoints).AddAttrib(gls.VertexPosition))
	pathMaterial := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	path := graphic.NewLineStrip(pathCircle, pathMaterial)
	moon := core.NewNode()
	moon.Add(distance)
	moon.Add(path)
	moon.RotateZ(LUNAR_PLANE_DEGREES * math32.Pi / 180)

	return moon
}
