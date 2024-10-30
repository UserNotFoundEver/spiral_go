package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math"
	"math/rand"
	"os"
	"time"
)

// Generate Fibonacci sequence for spiral radii with random start size
func fibonacciSequence(n int) []int {
	scale := rand.Intn(10) + 1 // Randomize the scale of the spiral
	sequence := make([]int, n)
	sequence[0], sequence[1] = scale, scale
	for i := 2; i < n; i++ {
		sequence[i] = sequence[i-1] + sequence[i-2]
	}
	return sequence
}

// Generate a psychedelic color pattern with randomness
func psychedelicColor(index, frame int) color.RGBA {
	r := uint8((index*15 + frame*5 + rand.Intn(50)) % 255)
	g := uint8((index*33 + frame*7 + rand.Intn(50)) % 255)
	b := uint8((index*55 + frame*11 + rand.Intn(50)) % 255)
	return color.RGBA{r, g, b, 255}
}

// Draw a gradient background with random shifts
func drawGradientBackground(img *image.RGBA, width, height int, frame int) {
	shiftR, shiftG, shiftB := rand.Intn(5), rand.Intn(5), rand.Intn(5) // Add randomness to gradient
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := uint8((x*shiftR + frame*5) % 255)
			g := uint8((y*shiftG + frame*7) % 255)
			b := uint8(((x + y) + frame*shiftB) % 255)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}
}

// Draw a Fibonacci spiral with random rotation and thickness
func drawFibonacciSpiral(img *image.RGBA, sequence []int, centerX, centerY int, frame int) {
	angleOffset := rand.Float64() * math.Pi / 4 // Randomize initial rotation
	angle := angleOffset
	for i, radius := range sequence {
		color := psychedelicColor(i, frame)

		// Define the bounding box for each arc in the spiral
		rect := image.Rect(
			centerX-int(radius), centerY-int(radius),
			centerX+int(radius), centerY+int(radius),
		)

		// Draw arc segment by rotating the start angle and thickening the line
		drawThickArc(img, rect, angle, angle+math.Pi/2, color, 3+rand.Intn(3)) // Random thickness
		angle += math.Pi / 2
	}
}

// Draws a thicker arc segment with a specified color and thickness
func drawThickArc(img *image.RGBA, rect image.Rectangle, startAngle, endAngle float64, clr color.Color, thickness int) {
	centerX := (rect.Min.X + rect.Max.X) / 2
	centerY := (rect.Min.Y + rect.Max.Y) / 2
	radius := (rect.Max.X - rect.Min.X) / 2

	for t := -thickness; t <= thickness; t++ { // Adds thickness by drawing multiple arcs
		for angle := startAngle; angle < endAngle; angle += 0.01 {
			x := centerX + int(float64(radius+t)*math.Cos(angle))
			y := centerY + int(float64(radius+t)*math.Sin(angle))
			img.Set(x, y, clr)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Seed for random data generation

	// Initialize the image dimensions
	width, height := 800, 800

	// Generate Fibonacci sequence for the spiral with random scaling
	sequence := fibonacciSequence(12) // Adjust for size of spiral
	centerX, centerY := width/2, height/2

	// Create animated GIF settings
	var images []*image.Paletted
	var delays []int

	// Generate frames
	numFrames := 30                     // Number of frames in the GIF
	delay := rand.Intn(10) + 5          // Random delay between frames for variety
	for frame := 0; frame < numFrames; frame++ {
		// Create a new image for each frame
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		// Draw gradient background for psychedelic effect with randomness
		drawGradientBackground(img, width, height, frame)

		// Draw the spiral with random color and rotation changes
		drawFibonacciSpiral(img, sequence, centerX, centerY, frame)

		// Convert RGBA image to Paletted image for GIF encoding
		palettedImg := image.NewPaletted(img.Bounds(), color.Palette{color.Black, color.White})
		draw.FloydSteinberg.Draw(palettedImg, img.Bounds(), img, image.Point{})

		// Append the frame and delay
		images = append(images, palettedImg)
		delays = append(delays, delay)
	}

	// Create a unique filename with a timestamp
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("unique_psychedelic_fibonacci_spiral_%s.gif", timestamp)
	outFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Save the GIF
	gif.EncodeAll(outFile, &gif.GIF{
		Image:     images,
		Delay:     delays,
		LoopCount: 0, // 0 = infinite loop
	})

	fmt.Printf("Unique Psychedelic Fibonacci Spiral GIF generated: %s\n", filename)
}
