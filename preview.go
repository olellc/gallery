package main

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"golang.org/x/image/draw"
)

// Preview returns image with the size of the longest dimension not greater then `size`.
// Supported src formats: jpeg, png, gif.
func Preview(src []byte, size int) ([]byte, error) {
	if size < 0 {
		size = 0
	}

	srcImg, _, err := image.Decode(bytes.NewReader(src))
	if err != nil {
		return nil, err
	}

	srcWidth := srcImg.Bounds().Dx()
	srcHeight := srcImg.Bounds().Dy()
	if srcWidth <= size && srcHeight <= size {
		return src, nil
	}

	var w, h int
	if srcWidth >= srcHeight {
		w, h = size, srcHeight*size/srcWidth
	} else {
		w, h = srcWidth*size/srcHeight, size
	}

	img := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.BiLinear.Scale(img, img.Bounds(), srcImg, srcImg.Bounds(), draw.Src, nil)

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
