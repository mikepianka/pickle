# pickle

Pickle is a little tool for picking pixels.

## Usage

Provide a JPG or PNG filepath and the X,Y coordinates of a pixel within the image. Pickle reads the image and can return the pixel value as either RGBA {uint8 uint8 uint8 uint8} or RGBa {uint8 uint8 uint8 float64}.

### Arguments

```bash
>>> pickle -h

Usage of ./pickle:
  -af
        Optional: Return alpha value as float instead of int (default false)
  -path string
        Path to JPG or PNG file
  -x int
        X coordinate of the pixel
  -y int
        Y coordinate of the pixel
```

### Examples

```bash
# color image with standard output
>>> pickle -path=./examples/gophers.jpg -x=300 -y=130

Pixel at (300, 130) has RGBA values of R: 90, G: 218, B: 255, A: 255
```

```bash
# color image with floating alpha output
>>> pickle -path=./examples/gophers.jpg -x=300 -y=130 -af

Pixel at (300, 130) has RGBA values of R: 90, G: 218, B: 255, A: 1.000000
```

```bash
# grayscale image with standard output
>>> pickle -path=./examples/gophers_grayscale.png -x=300 -y=130

Pixel at (300, 130) has RGBA values of R: 184, G: 184, B: 184, A: 255
```

```bash
# grayscale image with floating alpha output
>>> pickle -path=./examples/gophers_grayscale.png -x=300 -y=130 -af

Pixel at (300, 130) has RGBA values of R: 184, G: 184, B: 184, A: 1.000000
```
