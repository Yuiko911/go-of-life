package main

import (
	rd "math/rand"
)

func GenerateRandomField(field *[][]int) {
	for i, row := range *field {
		for j := range row {
			(*field)[i][j] = rd.Intn(2)
		}
	}
}

func GenerateFlyerField(field *[][]int) {
	(*field)[6][10] = 1
    (*field)[7][11] = 1
    (*field)[8][9] = 1
    (*field)[8][10] = 1
    (*field)[8][11] = 1
}

func ComputeNextField(field *[][]int) {

	y := len(*field)
	x := len((*field)[0])
	
	// Create a buffer field  
	var buffer[][] int = make([][]int, y)
	for i := range buffer {
		buffer[i] = make([]int, x)
	}

	// Counting the neighbours
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			buffer[i][j] = CountNeighbours(field, i, j)
		}
	}	
	
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			// ?
			if ((*field)[i][j] >= 1) {
				if (buffer[i][j] == 2 || buffer[i][j] == 3) {
					(*field)[i][j]++
				} else {
					(*field)[i][j] = 0
				}
			} else {
				if (buffer[i][j] == 3) {
					(*field)[i][j] = 1
				} else {
					(*field)[i][j] = 0
				}
			}

		}
	}
}

func CountNeighbours(field *[][]int, py int, px int) int {
	// Based on 
	// https://github.com/Yuiko911/silly/blob/main/game_of_life_clean.c
	
	y := len(*field)
	x := len((*field)[0])
	
	n := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i == 0 && j == 0) {continue}

			ny := (py + i + y) % y;
            nx := (px + j + x) % x;

			if (*field)[ny][nx] >= 1 {n++}
		}
	}

	return n;
}