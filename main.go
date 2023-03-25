package main

import (
	"flag"
	"fmt"
	"image"
	"log"
	"os"

	"image/color"
	_ "image/jpeg"
	_ "image/png"
)

type RGBa struct {
	R uint8
	G uint8
	B uint8
	A float64
}

type Config struct {
	imagePath string
	x         int
	y         int
	aFloat    bool
}

func main() {
	conf, err := cli()
	if err != nil {
		log.Fatal(err)
	}

	img, err := loadImage(conf.imagePath)
	if err != nil {
		log.Fatal(err)
	}

	rgba := getPixelValue(img, conf.x, conf.y)

	if conf.aFloat {
		a := alphaIntToFloat(rgba.A)
		fmt.Printf("Pixel at (%d, %d) has RGBA values of R: %d, G: %d, B: %d, A: %f\n", conf.x, conf.y, rgba.R, rgba.G, rgba.B, a)
	} else {
		fmt.Printf("Pixel at (%d, %d) has RGBA values of R: %d, G: %d, B: %d, A: %d\n", conf.x, conf.y, rgba.R, rgba.G, rgba.B, rgba.A)
	}
}

func cli() (Config, error) {
	imagePath := flag.String("path", "", "Path to JPG or PNG file")
	x := flag.Int("x", 0, "X coordinate of the pixel")
	y := flag.Int("y", 0, "Y coordinate of the pixel")
	aFloat := flag.Bool("af", false, "Optional: Return alpha value as float instead of int (default false)")

	flag.Parse()

	config := Config{
		imagePath: *imagePath,
		x:         *x,
		y:         *y,
		aFloat:    *aFloat,
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

func getPixelValue(img image.Image, x, y int) color.NRGBA {
	pixel := img.At(x, y)
	return color.NRGBAModel.Convert(pixel).(color.NRGBA)
}

func alphaIntToFloat(alpha uint8) float64 {
	return float64(alpha) / 255
}
