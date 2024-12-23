package main

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func main() {
	// List of hardcoded dimensions
	dimensions := []string{"size", "lookup", "update"}

	// List of directories as subjects
	dirs := os.Args[1:]
	if len(dirs) == 0 {
		fmt.Println("Usage: benchplot <dirs...>")
		os.Exit(1)
	}

	// Map to hold geomean values for each dimension and subject
	points := make(map[string]map[string]float64)

	// Process each directory and dimension
	for _, dir := range dirs {
		subject := filepath.Base(dir)
		points[subject] = make(map[string]float64)
		for _, dim := range dimensions {
			points[subject][dim] = extractGeomean(filepath.Join(dir, dim+".bm"))
		}
	}

	// Generate plots for each combination of dimensions
	for i := 0; i < len(dimensions); i++ {
		for j := i + 1; j < len(dimensions); j++ {
			dim1 := dimensions[i]
			dim2 := dimensions[j]
			createPlot(dim1, dim2, points)
		}
	}
}

// extractGeomean runs benchstat on a .bm file and extracts the geomean value
func extractGeomean(filepath string) float64 {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Printf("File not found: %s\n", filepath)
		return 0
	}

	// Run benchstat
	cmd := exec.Command("benchstat", filepath)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error running benchstat on %s: %v\n", filepath, err)
		os.Exit(1)
	}

	// Extract geomean from the output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "geomean") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				break
			}
			// Convert value to float
			value, err := parseValue(parts[len(parts)-1])
			if err == nil {
				return value
			}
		}
	}
	return 0
}

// parseValue parses a string like "59.30n" or "2.601Mi" into a float64
func parseValue(s string) (float64, error) {
	unitMultipliers := map[string]float64{
		"n":  1e-9,
		"u":  1e-6,
		"m":  1e-3,
		"K":  1e3,
		"Ki": 1024,
		"M":  1e6,
		"Mi": 1024 * 1024,
		"G":  1e9,
		"Gi": 1024 * 1024 * 1024,
	}

	// Split numeric part and unit
	for unit, multiplier := range unitMultipliers {
		if strings.HasSuffix(s, unit) {
			value, err := strconv.ParseFloat(strings.TrimSuffix(s, unit), 64)
			if err != nil {
				return 0, err
			}
			return value * multiplier, nil
		}
	}

	// Default to plain number
	return strconv.ParseFloat(s, 64)
}

func createPlot(dim1, dim2 string, points map[string]map[string]float64) {
	// Define units for each dimension
	units := map[string]string{
		"lookup": "ns",
		"size":   "KiB",
		"update": "ns",
	}

	// Validate that units exist for the given dimensions
	unit1, ok1 := units[dim1]
	unit2, ok2 := units[dim2]
	if !ok1 || !ok2 {
		fmt.Printf("Missing units for dimensions %s or %s\n", dim1, dim2)
		os.Exit(1)
	}

	p := plot.New()

	// Set plot title and axis labels
	p.Title.Text = fmt.Sprintf("%s vs %s", dim2, dim1)
	p.Title.TextStyle.Font.Size = vg.Points(16)
	p.X.Label.Text = fmt.Sprintf("%s (%s)", dim1, unit1)
	p.Y.Label.Text = fmt.Sprintf("%s (%s)", dim2, unit2)

	// Create data points with labels
	pts := make(plotter.XYs, len(points))
	labels := make([]string, len(points))
	idx := 0
	for subject, values := range points {
		// Convert units
		x := values[dim1]
		y := values[dim2]
		if unit1 == "ns" {
			x *= 1e9
		} else if unit1 == "KiB" {
			x /= 1024
		}
		if unit2 == "ns" {
			y *= 1e9
		} else if unit2 == "KiB" {
			y /= 1024
		}
		pts[idx].X = x
		pts[idx].Y = y
		labels[idx] = subject
		idx++
	}

	// Create scatter plot
	scatter, err := plotter.NewScatter(pts)
	if err != nil {
		panic(err)
	}
	p.Add(scatter)

	// Adjust axis ranges with padding
	xMin, xMax := math.Inf(1), math.Inf(-1)
	yMin, yMax := math.Inf(1), math.Inf(-1)
	for _, pt := range pts {
		if pt.X < xMin {
			xMin = pt.X
		}
		if pt.X > xMax {
			xMax = pt.X
		}
		if pt.Y < yMin {
			yMin = pt.Y
		}
		if pt.Y > yMax {
			yMax = pt.Y
		}
	}
	xPadding := (xMax - xMin) * 0.1
	yPadding := (yMax - yMin) * 0.1
	p.X.Min, p.X.Max = xMin-xPadding, xMax+xPadding
	p.Y.Min, p.Y.Max = yMin-yPadding, yMax+yPadding

	// Configure tick marks
	for _, axis := range []plot.Axis{p.X, p.Y} {
		axis.Tick.Marker = plot.TickerFunc(func(min, max float64) []plot.Tick {
			ticks := plot.DefaultTicks{}.Ticks(min, max)
			for i := range ticks {
				ticks[i].Label = fmt.Sprintf("%.0f", ticks[i].Value)
			}
			return ticks
		})
	}

	// Add labels to points
	for i, pt := range pts {
		label, err := plotter.NewLabels(plotter.XYLabels{
			XYs:    []plotter.XY{{X: pt.X + xMax*0.01, Y: pt.Y + yMax*0.01}},
			Labels: []string{labels[i]},
		})
		if err != nil {
			panic(err)
		}
		p.Add(label)
	}

	plotFilename := fmt.Sprintf("benchplot_%s_vs_%s.png", dim2, dim1)

	// Define the canvas size (with padding)
	canvasWidth := 8 * vg.Inch  // Total canvas width
	canvasHeight := 6 * vg.Inch // Total canvas height
	padding := 0.25 * vg.Inch   // Padding around the plot

	// Create the image canvas
	img := vgimg.New(canvasWidth, canvasHeight)
	dc := draw.New(img)

	// Create a rectangle defining the drawable area within the padded canvas
	cropped := draw.Canvas{
		Canvas: dc,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: padding, Y: padding},                              // Start after padding
			Max: vg.Point{X: canvasWidth - padding, Y: canvasHeight - padding}, // End before padding
		},
	}

	// Draw the plot in the cropped area
	p.Draw(cropped)

	// Save the image as a PNG
	w, err := os.Create(plotFilename)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	png := vgimg.PngCanvas{Canvas: img}
	_, err = png.WriteTo(w)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Plot saved as %s\n", plotFilename)
}
