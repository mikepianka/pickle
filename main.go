package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

type RGBA struct {
	r uint8
	g uint8
	b uint8
	a float64
}

type Config struct {
	imagePath string
	x         int
	y         int
}

func main() {
	config, err := cli()
	if err != nil {
		log.Fatal(err)
	}

	img, err := loadImage(config.imagePath)
	if err != nil {
		log.Fatal(err)
	}

	rgba := getPixelValue(img, config.x, config.y)

	fmt.Printf("Pixel at (%d, %d) has RGBA values of %v\n", config.x, config.y, rgba)
}

func cli() (Config, error) {
	imagePath := flag.String("path", "", "Path to JPG or PNG file")
	x := flag.Int("x", 0, "X coordinate of the pixel")
	y := flag.Int("y", 0, "Y coordinate of the pixel")

	flag.Parse()

	config := Config{
		imagePath: *imagePath,
		x:         *x,
		y:         *y,
	}

	if *imagePath == "" {
		return config, fmt.Errorf("invalid image path provided")
	}

	return config, nil
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func getPixelValue(img image.Image, x, y int) RGBA {
	pixel := img.At(x, y)
	r32, g32, b32, a32 := pixel.RGBA()

	// convert 16 bit values to 8 bit
	return RGBA{
		r: uint8(r32 / 257),
		g: uint8(g32 / 257),
		b: uint8(b32 / 257),
		a: float64(a32) / 65535.0,
	}
}
