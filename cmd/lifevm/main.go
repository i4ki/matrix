package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Matrix needs to replicate the state of early earth times
// It's not possible to do proper binary fission[1] of cell
// entities in modern operating systems and languages in
// a simpler way (it requires understanding of binary formats
// in the cell code)

// - Lifevm is an empty space
// - it runs primitive life forms
// - it simulates an aggressive environment (high temp., thunderstorms, etc)

const (
	SZ          = 8
	RefreshRate = 1 // in seconds
)

type (
	Pos uint8

	Point struct {
		x, y Pos
	}

	Plane struct {
		cells [SZ][SZ]Pos
	}

	LifeForm interface {
		Pos() Point
	}

	LifeVM struct {
		*sync.RWMutex
		space     Plane
		lifeforms []LifeForm
	}

	Archaea struct {
		pos Point
	}
)

func printSpace(p Plane) {
	for i := 0; i < SZ; i++ {
		for j := 0; j < SZ; j++ {
			if p.cells[i][j] == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("X")
			}
		}

		fmt.Printf("\n")
	}
}

func NewLifeVM() *LifeVM {
	vm := &LifeVM{
		RWMutex: &sync.RWMutex{},
	}

	go vm.updateSpace()
	return vm
}

func (vm *LifeVM) AddLife(life LifeForm) {
	vm.lifeforms = append(vm.lifeforms, life)
}

func (vm *LifeVM) updateSpace() {
	for {
		vm.Lock()
		for _, life := range vm.lifeforms {
			pos := life.Pos()

			for x := 0; x < SZ; x++ {
				for y := 0; y < SZ; y++ {
					if vm.space.cells[x][y] == 0 &&
						pos.x == Pos(x) && pos.y == Pos(y) {
						vm.space.cells[x][y]++
					}
				}
			}
		}

		time.Sleep(RefreshRate * time.Second)
		vm.Unlock()
	}
}

func NewArchaea() *Archaea {
	return &Archaea{
		pos: Point{
			Pos(rand.Intn(SZ)), Pos(rand.Intn(SZ)),
		},
	}
}

func (life *Archaea) Pos() Point {
	return life.pos
}

func main() {
	vm := NewLifeVM()
	life := NewArchaea()
	vm.AddLife(life)

	for {
		vm.RLock()
		printSpace(vm.space)
		vm.RUnlock()
		time.Sleep(1 * time.Second)
	}
}
