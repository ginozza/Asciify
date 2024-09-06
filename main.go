package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
)

const (
	width  = 100 
	height = 100 
	chars  = "@%#*+=-:. " 
)

func main() {
	// Define the flag for the image path
	imagePath := flag.String("image", "image.png", "Path to the image file")
	flag.Parse()

	// Open the image file
	file, err := os.Open(*imagePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Resize the image while maintaining aspect ratio
	img = resizeImage(img, width, height)

	// Convert the image to ASCII with color
	asciiArt := imageToASCIIWithColor(img)
	fmt.Println(asciiArt)
}

// Function to resize the image while maintaining aspect ratio
func resizeImage(img image.Image, maxWidth, maxHeight int) image.Image {
	bounds := img.Bounds()
	origWidth := bounds.Dx()
	origHeight := bounds.Dy()

	// Calculate the original aspect ratio
	aspectRatio := float64(origWidth) / float64(origHeight)

	// Define a factor to make the image wider
	const widthFactor = 1.2
	var newWidth, newHeight int

	// Adjust width and height to ensure the image is wider than tall
	if aspectRatio > 1 {
		newWidth = int(float64(maxWidth) * widthFactor)
		newHeight = int(float64(newWidth) / aspectRatio)
		if newHeight > maxHeight {
			newHeight = maxHeight
			newWidth = int(float64(newHeight) * aspectRatio)
		}
	} else {
		// Taller than wide or square image
		newHeight = int(float64(maxHeight) * widthFactor)
		newWidth = int(float64(newHeight) * aspectRatio)
		if newWidth > maxWidth {
			newWidth = maxWidth
			newHeight = int(float64(newWidth) / aspectRatio)
		}
	}

	// Ensure dimensions are not less than 1
	if newWidth <= 0 {
		newWidth = 1
	}
	if newHeight <= 0 {
		newHeight = 1
	}

	// Ensure width is greater than height
	if newWidth <= newHeight {
		newWidth = newHeight + 1
	}

	// Resize the image
	resizedImg := resize.Resize(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
	return resizedImg
}

// Function to convert the image to ASCII with color
func imageToASCIIWithColor(img image.Image) string {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	var asciiArt string

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Normalize color values
			r = r >> 8
			g = g >> 8
			b = b >> 8

			// Get the corresponding ASCII character
			gray := (r*299 + g*587 + b*114) / 1000
			index := int(gray * uint32(len(chars)) / 65536) // Correct the type

			if index >= len(chars) {
				index = len(chars) - 1
			}

			char := chars[index]
			// Add the character with color to the ASCII art
			asciiArt += colorize(char, int(r), int(g), int(b))
		}
		asciiArt += "\n"
	}

	return asciiArt
}

// Function to colorize the character
func colorize(char byte, r, g, b int) string {
	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", r, g, b, rune(char))
}
