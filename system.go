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

func revToSeconds(rx, ry float32) float32 {
	t := float32(0.0)

	if ry <= 0 && rx == 0 {
		t = (-1 * ry) * (15779076.48 / math32.Pi)
	} else if ry <= 0 && (rx < 0 || rx > 0) {
		t = (math32.Pi + ry) * (15779076.48 / math32.Pi)
	} else if ry >= 0 && rx < 0 {
		t = (math32.Pi + ry) * (15779076.48 / math32.Pi)
	} else if ry >= 0 && rx == 0 {
		t = (2*math32.Pi - ry) * (15779076.48 / math32.Pi)
	}

	return t
}

func main() {
	app := app.App()
	system := core.NewNode()
	gui.Manager().Set(system)

	cam := camera.New(1)
	cam.SetPosition(0, 7, 15)
	cam.LookAt(&math32.Vector3{X: 0.0, Y: 0.0, Z: 0.0}, &math32.Vector3{X: 0.0, Y: 10.0, Z: 10.0})
	camera.NewOrbitControl(cam)
	system.Add(cam)

	system.Add(newSun())
	earth := newEarth()
	system.Add(earth.planet)
	system.Add(earth.path)

	onResize := func(evname string, ev interface{}) {
		width, height := app.GetSize()
		app.Gls().Viewport(0, 0, int32(width), int32(height))
		cam.SetAspect(float32(width) / float32(height))
	}
	app.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	control := gui.NewHScrollBar(780, 20)
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

		timeInit := revToSeconds(earth.distance.Rotation().X,
			earth.distance.Rotation().Y)

		earth.planet.RotateY(delta)
		earth.distance.RotateY(-delta)
		earth.tilt.RotateY(delta * 365.2564)
		earth.moon.RotateY(delta * 365.2564 / 27.3)

		renderer.Render(system, cam)

		timeNew := revToSeconds(earth.distance.Rotation().X,
			earth.distance.Rotation().Y)

		dateTime := float32(0.0)
		if timeInit > 23668614 && timeNew < 7889538 {
			year += 31558152.96
		} else if timeInit < 7889538 && timeNew > 23668614 {
			year -= 31558152.96
		} else {
			dateTime = timeNew + year
		}

		dateTimeDisplay.SetText(fmt.Sprint(time.Unix(int64(dateTime), 0)))
	})
}
