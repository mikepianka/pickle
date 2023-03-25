# pickle

Pickle is a little tool for picking pixels.

## Usage

Provide a JPG or PNG filepath and the X,Y coordinates of a pixel within the image. Pickle reads the image and returns the pixel value as RGBA {uint8 uint8 uint8 float64}.

```bash
>>> pickle -h

Usage of ./pickle:
  -path string
        Path to JPG or PNG file
  -x int
        X coordinate of the pixel
  -y int
        Y coordinate of the pixel
```

```bash
>>> pickle -path=./examples/gophers.jpg -x=100 -y=200

Pixel at (100, 200) has RGBA values of {90 217 254 1}
```

```bash
>>> pickle -path=./examples/gophers_grayscale.png -x=100 -y=200

Pixel at (100, 200) has RGBA values of {184 184 184 1}
```
