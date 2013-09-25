package main

import (
	"flag"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"

	"code.google.com/p/graphics-go/graphics"
)

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	gifSeq, err := gif.DecodeAll(f)
	if err != nil {
		log.Fatal(err)
	}
	frameLength := len(gifSeq.Image)
	if frameLength < 2 {
		log.Fatal("Not enough frames")
	}
	bounds := gifSeq.Image[0].Bounds()
	bound := max(bounds.Max.X, bounds.Max.Y)
	// the distance between frames. The const bellow is random, change it
	// according what feels better.
	sliceLength := float64(bound) * 3.5
	cirleLength := sliceLength * float64(frameLength)

	// we are interested in the circle that is 2 times smaller to place the
	// pictures on it's "edges"
	radius := int(float64(cirleLength) / (2.0 * math.Pi))
	// not really a center..
	circleOrigin := &image.Point{radius, radius}

	slice := (2.0 * math.Pi) / float64(frameLength)
	outputImg := image.NewRGBA(image.Rect(0, 0, radius*2, radius*2))
	for i, img := range gifSeq.Image {
		// recalculate the angle
		angle := slice * float64(i+1)
		x := circleOrigin.X + int(float64(radius/2)*math.Cos(angle))
		y := circleOrigin.Y + int(float64(radius/2)*math.Sin(angle))
		halfX := bounds.Max.X
		halfY := bounds.Max.Y
		minX, minY, maxX, maxY := x-halfX, y-halfY, x+halfX, y+halfY
		// make a holder for a new image
		rotatedImg := image.NewRGBA(image.Rect(0, 0, bound, bound))
		// rotate the image using the external lib
		graphics.Rotate(rotatedImg, img, &graphics.RotateOptions{angle + math.Pi/2})
		// draw it in the global image
		draw.Draw(outputImg, image.Rect(minX, minY, maxX, maxY), rotatedImg, image.ZP, draw.Src)
	}
	outputFp, err := os.Create("output.png")
	png.Encode(outputFp, outputImg)
}
