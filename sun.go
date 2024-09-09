package main

import (
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
)

func createSun() *graphic.Mesh {
	sunShape := geometry.NewSphere(2, 360, 360)
	sunTexture := material.NewStandard(&math32.Color{R: 1.0, G: 0.8, B: 0.5})
	sunTexture.SetEmissiveColor(&math32.Color{R: 1.0, G: 0.8, B: 0.5})
	sunLight := light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 40.0)
	sunLight.SetPosition(0.0, 0.0, 0.0)
	sun := graphic.NewMesh(sunShape, sunTexture)
	sun.Add(light.NewAmbient(&math32.Color{R: 1.0, G: 1.0, B: 1.0}, 0.2))
	sun.Add(sunLight)

	return sun
}
