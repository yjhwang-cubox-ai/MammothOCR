package utils

import (
	"bytes"
	"image"
)

func BytesToImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

func CropImage(img image.Image, box BoundingBox) image.Image {
	bounds := img.Bounds()
	x1 := int(box.X1 * float32(bounds.Max.X))
	y1 := int(box.Y1 * float32(bounds.Max.Y))
	x2 := int(box.X2 * float32(bounds.Max.X))
	y2 := int(box.Y2 * float32(bounds.Max.Y))

	return img.(*image.RGBA).SubImage(image.Rect(x1, y1, x2, y2))
}
