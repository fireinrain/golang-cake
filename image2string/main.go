package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"os"
)

func main() {
	// Open the image file
	file, err := os.Open("example.jpg")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the image file
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// Resize the image to a smaller size for better results
	width, height := img.Bounds().Max.X, img.Bounds().Max.Y
	aspectRatio := float64(height) / float64(width)
	newWidth := 100
	newHeight := int(float64(newWidth) * aspectRatio)
	img = resize(img, newWidth, newHeight)

	// Convert the image to grayscale
	img = grayscale(img)

	// Define the characters to use for the ASCII art
	chars := []string{" ", ".", "*", ":", "o", "&", "8", "#", "@"}

	// Convert the image to a string of characters
	charPixels := convertToCharPixels(img, chars)

	// Print the resulting ASCII art
	fmt.Println(charPixels)
}

// resize resizes an image to the specified width and height
func resize(img image.Image, width, height int) image.Image {
	resized := image.NewRGBA(image.Rect(0, 0, width, height))
	//draw.NearestNeighbor.Resample(resized, img.Bounds(), img, img.Bounds().Min, draw.Over, nil)
	return resized
}

// grayscale converts an image to grayscale
func grayscale(img image.Image) image.Image {
	gray := image.NewGray(img.Bounds())
	draw.Draw(gray, gray.Bounds(), img, img.Bounds().Min, draw.Over)
	return gray
}

// convertToCharPixels converts an image to a string of characters
func convertToCharPixels(img image.Image, chars []string) string {
	pixels := img.Bounds().Max.X * img.Bounds().Max.Y
	charPixels := make([]string, pixels)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			pixel := img.At(x, y)
			grayPixel := color.GrayModel.Convert(pixel).(color.Gray)
			index := int(grayPixel.Y * 9 / 255)
			charPixels[x+y*img.Bounds().Max.X] = chars[index]
		}
	}
	return fmt.Sprintf("%s", charPixels)
}
