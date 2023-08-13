package helper

import (
	"image"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
)

func GenerateImageWithWrappedText(fontFile string, text string, outputImageFile string) error {
	const canvasWidth = 800
	const canvasHeight = 200
	const fontSize = 48

	dc := gg.NewContext(canvasWidth, canvasHeight)

	dc.SetRGB(0, 0, 0)

	dc.Clear()

	if fontFile != "" {
		err := dc.LoadFontFace(fontFile, fontSize)
		if err != nil {
			log.Fatal(err)
		}
	}

	dc.SetRGB(255, 255, 255) // Set text color to black
	dc.SetRGB255(255, 255, 255)
	dc.DrawStringWrapped(text, 0, 0, 0, 0, canvasWidth, 1.5, gg.AlignCenter)

	err := dc.SavePNG(outputImageFile)
	if err != nil {
		log.Fatal(err)
	}

	// Open the generated image file
	file, err := os.Open(outputImageFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Get the bounds of the image
	bounds := img.Bounds()

	// Find the minimum and maximum x and y coordinates of the non-black pixels
	minX := bounds.Max.X
	maxX := 0
	minY := bounds.Max.Y
	maxY := 0

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if r != 0 || g != 0 || b != 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	// Crop the image using the calculated coordinates
	croppedImg := img.(interface {
		SubImage(r image.Rectangle) image.Image
	}).SubImage(image.Rect(minX, minY, maxX+1, maxY+1))

	// Save the cropped image to a new file
	croppedFile, err := os.Create(outputImageFile)
	if err != nil {
		return err
	}
	defer croppedFile.Close()

	err = png.Encode(croppedFile, croppedImg)
	if err != nil {
		return err
	}

	return nil
}
