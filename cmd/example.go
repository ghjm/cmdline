package main

import (
	"fmt"
	"github.com/ghjm/cmdline"
	"os"
)

type circle struct {
	Radius float64 `description:"Radius of the circle" barevalue:"yes" required:"yes"`
	Color  string  `description:"Color to use when drawing the circle" default:"white"`
}

func (c circle) Check() error {
	if c.Radius < 0 {
		return fmt.Errorf("circle radius cannot be negative")
	}
	return nil
}

func (c circle) Draw() error {
	fmt.Printf("Drawing a %s circle of radius %f\n", c.Color, c.Radius)
	return nil
}

type rectangle struct {
	Width  float64 `description:"Width of the rectangle" required:"yes"`
	Height float64 `description:"Height of the rectangle" required:"yes"`
	Color  string  `description:"Color to use when drawing the rectangle" default:"white"`
}

func (r rectangle) Check() error {
	if r.Height < 0 || r.Width < 0 {
		return fmt.Errorf("rectangle height and width cannot be negative")
	}
	return nil
}

func (r rectangle) Draw() error {
	fmt.Printf("Drawing a %s rectangle of height %f and width %f\n", r.Color, r.Height, r.Width)
	return nil
}

func main() {
	cl := cmdline.NewCmdline()
	cl.AddConfigType("circle", "Circle Shape", circle{})
	cl.AddConfigType("rectangle", "Rectangle Shape", rectangle{})
	_, err := cl.ParseAndRun(os.Args[1:], []string{"Check", "Draw"}, cmdline.ShowHelpIfNoArgs)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
