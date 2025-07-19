package main

import (
	"flag"
	"log"
	"slices"

	gc "github.com/rthornton128/goncurses"
)

type FieldStatus struct {
	Monochrome bool
	LightMode  bool

	Paused       bool
	CurrentSpeed int
}

func main() {

	fs := FieldStatus{
		Monochrome: false,
		LightMode:  false,

		Paused:       false,
		CurrentSpeed: 50,
	}

	flag.BoolVar(&fs.Paused, "p", false, "Start with the simulation paused.")
	flag.BoolVar(&fs.Monochrome, "m", false, "Start in monochrome mode.")
	flag.BoolVar(&fs.LightMode, "l", false, "Enable light mode.")

	flagSpeed := flag.Int64("i", 50, "Initial speed of the simulation.")

	flag.Parse() // This line *has* to be before gc.Init()

	fs.CurrentSpeed = int(*flagSpeed)

	scr, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	if err = gc.StartColor(); err != nil {
		log.Fatal("Colors are not supported on this terminal")
	}

	gc.Cursor(0)
	gc.Echo(false)
	gc.Raw(true)

	gc.UseDefaultColors()
	SetColors(fs)

	wy, wx := scr.MaxYX()
	var field [][]int = make([][]int, wy)
	for i := range field {
		field[i] = make([]int, wx)
	}

	scr.Timeout(fs.CurrentSpeed)
	GenerateRandomField(&field)

	for {
		DrawToScreen(scr, fs, field, 0, 0)
		scr.Refresh()

		if !fs.Paused {
			ComputeNextField(&field)
		}

		ch := scr.GetChar()
		switch ch {
		case 'q':
			return
		case 'r':
			GenerateRandomField(&field)
		case 's':
			fs.ChangeSpeed()
			scr.Timeout(fs.CurrentSpeed)
		case 'c':
			fs.ToggleColors()
		case 'p':
			fs.TogglePause()
		}
	}
}

func (fs *FieldStatus) ToggleColors() {
	fs.Monochrome = !fs.Monochrome
	SetColors(*fs)
}

func (fs *FieldStatus) ChangeSpeed() {
	speeds := []int{10, 25, 50, 100, 200, 500, 1000}

	cs := slices.Index(speeds, fs.CurrentSpeed)

	// In case fs.CurrentSpeed is *somehow* not in the array,
	// cs will be -1, which will put the speed at speeds[0].

	if cs >= len(speeds)-1 {
		fs.CurrentSpeed = speeds[0]
	} else {
		fs.CurrentSpeed = speeds[cs+1]
	}

}

func (fs *FieldStatus) TogglePause() {
	fs.Paused = !fs.Paused
}

func SetColors(fs FieldStatus) {
	if fs.Monochrome {
		gc.InitPair(1, -1, -1)
		gc.InitPair(2, -1, -1)
		var i int16
		for i = 3; i < 10; i++ {
			gc.InitPair(i, -1, gc.C_WHITE)
		}

	} else if fs.LightMode {
		gc.InitPair(1, -1, -1)
		gc.InitPair(2, gc.C_BLACK, gc.C_WHITE)
		gc.InitPair(3, gc.C_BLACK, gc.C_BLACK)
		gc.InitPair(4, gc.C_BLACK, gc.C_CYAN)
		gc.InitPair(5, gc.C_WHITE, gc.C_BLUE)
		gc.InitPair(6, gc.C_BLACK, gc.C_GREEN)
		gc.InitPair(7, gc.C_BLACK, gc.C_YELLOW)
		gc.InitPair(8, gc.C_BLACK, gc.C_MAGENTA)
		gc.InitPair(9, gc.C_BLACK, gc.C_RED)
	} else {
		gc.InitPair(1, -1, -1)
		gc.InitPair(2, gc.C_WHITE, gc.C_BLACK)
		gc.InitPair(3, gc.C_WHITE, gc.C_WHITE)
		gc.InitPair(4, gc.C_WHITE, gc.C_RED)
		gc.InitPair(5, gc.C_BLACK, gc.C_MAGENTA)
		gc.InitPair(6, gc.C_WHITE, gc.C_YELLOW)
		gc.InitPair(7, gc.C_WHITE, gc.C_GREEN)
		gc.InitPair(8, gc.C_WHITE, gc.C_BLUE)
		gc.InitPair(9, gc.C_WHITE, gc.C_CYAN)
	}
}

func DrawToScreen(scr *gc.Window, fs FieldStatus, field [][]int, y int, x int) {
	var c int16 = 1

	for i := range field {
		for j := range field[i] {

			if field[i][j] <= 7 {
				c = (int16)(field[i][j]) + 2
			} else {
				c = 9
			}

			scr.ColorOn(c)
			scr.MovePrint(i+y, j+x, " ")
			scr.ColorOff(c)

		}
	}
}
