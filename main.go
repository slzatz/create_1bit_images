package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dropbox/godropbox/container/bitvector"
	"github.com/slzatz/images/bmp1bit"
)

func main() {

	fname := os.Args[1]
	name, _, _ := strings.Cut(fname, ".")
	cmd := exec.Command(
		"convert", fname, "-resize", "296x400^", "-gravity", "North",
		"-ordered-dither", "o8x8", "-extent", "296x400",
		"-depth", "1", "-type", "bilevel", "BMP3:"+name+"8d.bmp")
	cmd.Dir = "/home/slzatz/code/images/pictures"
	err := cmd.Run()
	if err != nil {
		log.Fatalf("ImageMagick convert error: %v\n", err)
	}

	file, _ := os.Open("pictures/" + name + "8d.bmp")
	defer file.Close()

	img, err := bmp1bit.Decode(file)
	if err != nil {
		log.Fatalf("Decode error: %v\n", err)
	}

	z, ok := img.(*image.Paletted)
	if !ok {
		log.Fatalf("Image is not Paletted")
	}

	fmt.Printf("Paletted\n")
	fmt.Printf("Number of pixels = %d\n", len(z.Pix))
	fmt.Printf("Stride = %v\n", z.Stride)
	fmt.Printf("Rect Min X = %v\n", z.Rect.Min.X)
	fmt.Printf("Rect Min Y = %v\n", z.Rect.Min.Y)
	fmt.Printf("Rect Max X = %v\n", z.Rect.Max.X)
	fmt.Printf("Rect Max Y = %v\n", z.Rect.Max.Y)
	bv := bitvector.NewBitVector([]byte{}, 0)
	for _, v := range z.Pix {
		bv.Append(v)
	}
	os.WriteFile("pictures/"+name+".bd", bv.Bytes(), 0644)
}
