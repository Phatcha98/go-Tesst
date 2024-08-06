package main

import (
	"gocv.io/x/gocv"
)

func main() {
	// Open webcam
	webcam, err := gocv.OpenVideoCapture(2)
	if err != nil {
		println("Error opening webcam:", err)
		return
	}
	defer webcam.Close()

	// Create window
	window := gocv.NewWindow("Camera Feed")
	defer window.Close()

	// Main loop to read and display frames
	for {
		// Read frame from webcam
		frame := gocv.NewMat()
		if ok := webcam.Read(&frame); !ok {
			println("Error reading frame from webcam")
			return
		}

		// Check if frame is empty
		if frame.Empty() {
			println("Empty frame from webcam")
			continue
		}

		// Display frame in window
		window.IMShow(frame)

		// Wait for key press for 1 millisecond
		key := window.WaitKey(1)

		// Check if the escape key was pressed
		if key == 27 {
			break
		}

		// Close frame
		frame.Close()
	}
}
