package main

import (
	"flag"
	"log"
	"slices"

	gc "github.com/rthornton128/goncurses"
)

type FieldStatus struct {
	Fullscreen bool
	Unicode    bool
	LightMode  bool

	Monochrome   bool
	Paused       bool
	CurrentSpeed int
}

func main() {
	/*
		-f --fullscreen
		-p --paused
		-i --initial-speed=[speed]

		-u --unicode

		-m --monochrome
		-L --light-mode
		-D --dark-mode

		-s --starting-state=[state]
			0: random
			1: checkered
	*/

	fs := FieldStatus{
		Fullscreen: false,
		Unicode:    false,
		LightMode:  false,

		Monochrome:   false,
		Paused:       false,
		CurrentSpeed: 50,
	}

	flag.BoolVar(&fs.Fullscreen, "f", false, "Start in fullscreen.") // Replace with size / padding
	flag.BoolVar(&fs.Paused, "p", false, "Start with the simulation paused.")

	flagSpeed := flag.Int64("i", 50, "Initial speed of the simulation.")
	fs.CurrentSpeed = int(*flagSpeed)

	flag.BoolVar(&fs.Monochrome, "m", false, "Start in monochrome mode.")
	flag.BoolVar(&fs.Unicode, "U", false, "Start in unicode mode (only use character). Force the -m flag.")

	flag.Parse() // This line *has* to be before gc.Init()

	scr, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	colorAvailable := true

	if err = gc.StartColor(); err != nil {
		log.Print("Colors are not supported on this terminal")
		colorAvailable = false
		fs.Monochrome = true
	}

	if fs.Unicode {
		colorAvailable = false
		fs.Monochrome = true
	}
	
	gc.Cursor(0)
	gc.Echo(false)
	gc.Raw(true)
	
	log.Print(fs.Monochrome)

	if colorAvailable {
		gc.UseDefaultColors()
		SetColors(fs)
	}

	wy, wx := scr.MaxYX()
	if wy < 20 || wx < 64 {
		log.Fatal("Terminal is not big enough")
	}

	wy -= 10 // Magic numbers for padding
	wx -= 6
	var field [][]int = make([][]int, wy)
	for i := range field {
		field[i] = make([]int, wx)
	}

	GenerateRandomField(&field)
	// GenerateCheckerboardField(&field)
	// GenerateFlyerField(&field)

	DrawBorderAroundField(scr, 3, 3, wy, wx)

	scr.Timeout(fs.CurrentSpeed)

	for {
		DrawMenu(scr, 3+wy+2, fs)
		DrawToScreen(scr, field, 3, 3)
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

func DrawToScreen(scr *gc.Window, field [][]int, y int, x int) {

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

func DrawBorderAroundField(scr *gc.Window, fy int, fx int, fsy int, fsx int) {
	scr.ColorOn(3)

	for i := 0; i < fsy+2; i++ {
		scr.MovePrint(fy-1+i, fx-1, " ")
		scr.MovePrint(fy-1+i, fx+fsx, " ")
	}

	for i := 0; i < fsx; i++ {
		scr.MovePrint(fy-1, fx+i, " ")
		scr.MovePrint(fy+fsy, fx+i, " ")
	}

	scr.ColorOff(3)
}

func DrawMenu(scr *gc.Window, y int, fs FieldStatus) {
	menu := "(q)uit     (r)efresh     (s)peed     toggle (c)olors     (p)ause"
	meop := "                         %-7d     %-15s"

	var colorstatus string
	if fs.Monochrome {
		colorstatus = "monochrome"
	} else {
		colorstatus = "colorful"
	}

	_, wx := scr.MaxYX()

	scr.MovePrint(y, (wx/2)-(len(menu)/2), menu)
	scr.MovePrintf(y+1, (wx/2)-(len(menu)/2), meop, fs.CurrentSpeed, colorstatus)
}

func SetColors(fs FieldStatus) {
	if fs.Monochrome {
		gc.InitPair(1, -1, -1)
		gc.InitPair(2, -1, -1)
		var i int16
		for i = 3; i < 10; i++ {
			gc.InitPair(i, -1, gc.C_WHITE)
		}
	} else {
		gc.InitPair(1, -1, -1)
		gc.InitPair(2, gc.C_WHITE, gc.C_BLACK)
		gc.InitPair(3, gc.C_WHITE, gc.C_WHITE)
		gc.InitPair(4, gc.C_WHITE, gc.C_CYAN)
		gc.InitPair(5, gc.C_BLACK, gc.C_BLUE)
		gc.InitPair(6, gc.C_WHITE, gc.C_GREEN)
		gc.InitPair(7, gc.C_WHITE, gc.C_YELLOW)
		gc.InitPair(8, gc.C_WHITE, gc.C_MAGENTA)
		gc.InitPair(9, gc.C_WHITE, gc.C_RED)

	}
}
