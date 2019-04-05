package main

import (
	"flag"
	"log"
	"image"
	"golang.org/x/image/math/fixed"
	"os"
	"image/png"
	"github.com/x6uderek/expgraphy/exp"
)

var input = flag.String("exp", "", "math expression")
var output = flag.String("out", "./out.png", "output file")
var fontSize = flag.Float64("font", 36.0, "font size in float64")

func main() {
	flag.Parse()
	if *input == "" {
		log.Fatal("empty input")
	}
	expr, err := exp.Parse(*input)
	if err != nil {
		log.Fatalf("parse err: %v",err)
	}
	gra := expr.GetGraphy()
	size := gra.Measure(*fontSize)
	img := image.NewRGBA(image.Rect(0,0, int(exp.FixedToFloat(size.X)*2), int(exp.FixedToFloat(size.Y)*2)))
	startPoint := fixed.Point26_6{X:size.X/2, Y:size.Y*3/2}
	gra.Draw(img, startPoint)

	f,err := os.Create(*output)
	if err!=nil {
		log.Fatalf("output file: %v", err)
	}
	defer f.Close()
	png.Encode(f, img)
}
