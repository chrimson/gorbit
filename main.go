package main

import (
	"fmt"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/window"

	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
)

const Q1_SECONDS = 7889538.24
const Q2_SECONDS = 15779076.48
const Q3_SECONDS = 23668614.72
const YEAR_SECONDS = 31558152.96
const REVOLUTION_DAYS = 365.2564
const LUNAR_DAYS = 27.3
const LUNAR_PLANE_DEGREES = 5.14
const EARTH_TILT_DEGREES = 23.4
const SUMMER_SOLSTICE_OFFSET = 1718841600.0 // 2024-06-20

func revToSeconds(rotationX, rotationY float32) (time, realY float32) {
	factor := float32(0.0)

	if rotationY <= 0 {
		if rotationX == 0 {
			factor = -rotationY
		} else {
			factor = math32.Pi + rotationY
		}
	} else {
		if rotationX == 0 {
			factor = 2*math32.Pi - rotationY
		} else {
			factor = math32.Pi + rotationY
		}
	}

	time = factor * float32(Q2_SECONDS) / math32.Pi
	realY = -factor

	return time, realY
}

func main() {
	app := app.App(1600, 800, "Gorbit")
	system := core.NewNode()
	gui.Manager().Set(system)

	cam := camera.New(1)
	cam.SetPosition(0.0, 7.0, 15.0)
	cam.LookAt(&math32.Vector3{X: 0.0, Y: 0.0, Z: 0.0}, &math32.Vector3{X: 0.0, Y: 1.0, Z: 0.0})
	c := camera.NewOrbitControl(cam)
	system.Add(cam)

	system.Add(newSun())
	earth := newEarth()
	system.Add(earth.body)
	system.Add(earth.path)

	onResize := func(evname string, ev interface{}) {
		width, height := app.GetSize()
		app.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	}
	app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	radioSun := gui.NewRadioButton("Sun")
	radioSun.SetPosition(10, 60)
	radioSun.SetGroup("center_body")
	radioSun.Subscribe(gui.OnChange, func(name string, ev interface{}) {
		if radioSun.Value() {
			cam.LookAt(&math32.Vector3{X: 0.0, Y: 0.0, Z: 0.0}, &math32.Vector3{X: 0.0, Y: 1.0, Z: 0.0})
			c.SetTarget(math32.Vector3{X: 0.0, Y: 0.0, Z: 0.0})
		}
	})
	system.Add(radioSun)

	radioEarth := gui.NewRadioButton("Earth")
	radioEarth.SetPosition(80, 60)
	radioEarth.SetGroup("center_body")
	radioEarth.Subscribe(gui.OnChange, func(name string, ev interface{}) {
	})
	system.Add(radioEarth)

	control := gui.NewHScrollBar(1580, 20)
	control.SetColor(&math32.Color{R: 0.2, G: 0.2, B: 0.2})
	control.SetPosition(10, 10)
	control.SetValue(0.50)
	control.Subscribe(gui.OnChange, func(evname string, ev interface{}) {})
	system.Add(control)

	dateTimeDisplay := gui.NewLabel("")
	dateTimeDisplay.SetPosition(10, 40)
	dateTimeDisplay.SetColor(&math32.Color{R: 1.0, G: 1.0, B: 1.0})
	system.Add(dateTimeDisplay)

	year := float32(0.0)
	app.Run(func(renderer *renderer.Renderer, deltaTime time.Duration) {
		app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

		speed := 0.0
		if control.Value() < 0.475 {
			speed = control.Value() - 0.475
		} else if control.Value() > 0.525 {
			speed = control.Value() - 0.525
		}
		delta := float32(speed) * float32(deltaTime.Seconds())

		timeInit, realY := revToSeconds(earth.distance.Rotation().X, earth.distance.Rotation().Y)

		earth.body.RotateY(delta)
		earth.distance.RotateY(-delta)
		earth.tilt.RotateY(delta + delta*REVOLUTION_DAYS)
		earth.moon.RotateY(delta * REVOLUTION_DAYS / LUNAR_DAYS)

		if radioEarth.Value() {
			cam.LookAt(&math32.Vector3{X: 10 * math32.Cos(realY),
				Y: 0.0,
				Z: 10 * math32.Sin(realY)},
				&math32.Vector3{X: 0.0, Y: 1.0, Z: 0.0})
			c.SetTarget(math32.Vector3{X: 10 * math32.Cos(realY),
				Y: 0.0,
				Z: 10 * math32.Sin(realY)})
		}

		app.Renderer().Render(system, cam)

		timeNew, _ := revToSeconds(earth.distance.Rotation().X, earth.distance.Rotation().Y)

		dateTime := float32(0.0)
		if timeInit > Q3_SECONDS && timeNew < Q1_SECONDS {
			year += YEAR_SECONDS
		} else if timeInit < Q1_SECONDS && timeNew > Q3_SECONDS {
			year -= YEAR_SECONDS
		} else {
			dateTime = timeNew + year + SUMMER_SOLSTICE_OFFSET
		}

		dateTimeDisplay.SetText(fmt.Sprint(time.Unix(int64(dateTime), 0).UTC()))
	})
}
