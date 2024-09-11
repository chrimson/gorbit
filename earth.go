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

type Earth struct {
	body     *core.Node
	distance *core.Node
	path     *graphic.LineStrip
	tilt     *graphic.Mesh
	moon     *core.Node
}

func newEarth() Earth {
	earth := Earth{}

	shape := geometry.NewSphere(0.5, 360, 360)
	image := func(path string) *texture.Texture2D {
		t, _ := texture.NewTexture2DFromImage(path)
		t.SetFlipY(false)
		return t
	}
	texture := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	texture.SetShininess(5)
	texture.AddTexture(image("earth.jpg"))
	tilt := graphic.NewMesh(shape, texture)
	tilt.RotateZ(EARTH_TILT_DEGREES * math32.Pi / 180)

	axisGeometry := geometry.NewGeometry()
	axisVertices := math32.NewArrayF32(0, 0)
	axisVertices.Append(
		0.0, 1.0, 0.0,
		0.0, -1.0, 0.0,
	)
	axisGeometry.AddVBO(gls.NewVBO(axisVertices).AddAttrib(gls.VertexPosition))
	axisMaterial := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	axis := graphic.NewLines(axisGeometry, axisMaterial)
	tilt.Add(axis)

	moon := newMoon()

	distance := core.NewNode()
	distance.Add(moon)
	distance.Add(tilt)
	distance.TranslateX(10.0)

	body := core.NewNode()
	body.Add(distance)

	pathCircle := geometry.NewGeometry()
	pathPoints := math32.NewArrayF32(0, 0)
	for x := float32(-1.0); x < 1.0; x = x + 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		pathPoints.Append(10.0*x, 0.0, 10.0*z)
	}
	for x := float32(1.0); x > -1.0; x = x - 0.01 {
		z := math32.Sqrt(1.0 - math32.Pow(x, 2))
		pathPoints.Append(10.0*x, 0.0, -10.0*z)
	}
	pathCircle.AddVBO(gls.NewVBO(pathPoints).AddAttrib(gls.VertexPosition))
	pathMaterial := material.NewStandard(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	path := graphic.NewLineStrip(pathCircle, pathMaterial)

	earth.body = body
	earth.distance = distance
	earth.path = path
	earth.tilt = tilt
	earth.moon = moon

	return earth
}
