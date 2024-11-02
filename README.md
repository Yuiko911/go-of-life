# Go of Life 🦠
Conway's Game of Life, written in Go using ncurses.

![alt test](readme/resized.gif)

## Dependencies
[rthornton128/goncurses](https://github.com/rthornton128/goncurses)

## Build and run
```sh
go build -o gooflife *.go
./gooflife
```
Alternatively, run without building:
```sh
go run *.go
```

## Usage
In the program :
- `s` changes the speed of the simulation
- `c` toggle on and off colors
- `p` pauses the simulation
- `r` regenerate a new field
- `q` quits
