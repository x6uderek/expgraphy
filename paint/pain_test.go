package paint

import (
	"testing"
	"image"
	"image/color"
	"os"
	"image/png"
)

func TestDrawLine(t *testing.T) {
	img := image.NewRGBA(image.Rect(0,0, 400, 400))
	DrawLine(img, 20, 380, 200, 20, 5, color.Black)
	DrawLine(img, 20,30, 350,300, 2, color.RGBA{230,240,0, 255})
	DrawLine(img, 375, 40, 200, 400, 3, color.RGBA{250, 0 ,180, 255})
	DrawLine(img, 380,380,0,20, 3, color.RGBA{0,250, 90, 255})
	DrawLine(img, 1,1, 399,1,2,color.RGBA{0xaa, 0x05, 0xcc, 255})
	DrawLine(img, 1,1,1,399, 2,color.RGBA{0x22,0xcc,0xff,255})
	f,err := os.Create(`F:\gb\bin\out1.png`)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}
