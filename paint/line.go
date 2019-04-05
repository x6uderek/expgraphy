package paint

import (
	"image/draw"
	"math"
	"image/color"
)

func DrawLine(image draw.Image, fx,fy, tx,ty, width float64, col color.Color) {
	dx := tx - fx
	dy := ty - fy
	d := math.Hypot(dx, dy)
	if d == 0 {
		return
	}
	sin := dy/d
	cos := dx/d
	x0 := fx - sin * width/2
	y0 := fy + cos * width/2
	for fDistance := 0.0; fDistance<=width; fDistance+=0.5 {
		startX := x0 + fDistance*sin
		startY := y0 - fDistance*cos

		for distance := 0.0; distance<=d; distance+=0.5 {
			image.Set(int(startX + distance*cos), int(startY + distance*sin), col)
		}
	}
}