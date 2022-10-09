package main

import (
	"fmt"
	"math/rand"
	"sync"
)

const I = 10
const N = 5
var matrix [N][N]float64
var cache [N][N]float64

type Item struct {
	i int
	j int
	v float64
}

func printMatrix() {
	fmt.Println()
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("[%v, %v] = %v\n", i, j, matrix[i][j])
		}
	}
}

func fillMatrix() {
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			x := float64(rand.Int())
			if x < 0 { continue }
			matrix[i][j] = x
		}
	}
}

func calculateNorthIndex(i int, j int) (int, int) {
	return i - 1, j
}

func calculateSouthIndex(i int, j int) (int, int) {
	return i + 1, j
}

func calculateEastIndex(i int, j int) (int, int) {
	return i, j - 1
}

func calculateWestIndex(i int, j int) (int, int) {
	return i, j + 1
}

func calculateAvg(i int, j int) float64 {
	sum := 0.0

	var coordI int
	var coordJ int
	var coordsI [4]int
	var coordsJ [4]int

	coordI, coordJ = calculateNorthIndex(i, j)
	if(coordI < 0 || coordJ < 0 || coordI >= N || coordJ >= N) { return -1 }
	coordsI[0] = coordI
	coordsJ[0] = coordJ

	coordI, coordJ = calculateSouthIndex(i, j)
	if(coordI < 0 || coordJ < 0 || coordI >= N || coordJ >= N) { return -1 }
	coordsI[1] = coordI
	coordsJ[1] = coordJ

	coordI, coordJ = calculateEastIndex(i, j)
	if(coordI < 0 || coordJ < 0 || coordI >= N || coordJ >= N) { return -1 }
	coordsI[2] = coordI
	coordsJ[2] = coordJ

	coordI, coordJ = calculateWestIndex(i, j)
	if(coordI < 0 || coordJ < 0 || coordI >= N || coordJ >= N) { return -1 }
	coordsI[3] = coordI
	coordsJ[3] = coordJ

	for i := 0; i < 4; i++ {
		sum += matrix[coordsI[i]][coordsJ[i]]
	}

	return sum / 4.0
}

func main() {
	fillMatrix()
	printMatrix()

	for i := 0; i < I; i++ {
		avgChannel := make(chan Item, 50)
		avgWaitGroup := sync.WaitGroup{}
		avgWaitGroup.Add(N * N)

		go func(ch <-chan Item, wg *sync.WaitGroup, cache *[N][N]float64) {
			for x := range ch {
				cache[x.i][x.j] = x.v
				wg.Done()
			}
		}(avgChannel, &avgWaitGroup, &cache)

		go func(ch chan<- Item) {
			for i := 0; i < N; i++ {
				for j := 0; j < N; j++ {
					x := calculateAvg(i, j)
					ch <- Item{i, j, x}
				}
			}
			close(ch)
		}(avgChannel)

		avgWaitGroup.Wait()

		writeChannel := make(chan Item, 50)
		writeWaitGroup := sync.WaitGroup{}
		writeWaitGroup.Add(N * N)

		go func(ch <-chan Item, wg *sync.WaitGroup, matrix *[N][N]float64) {
			for x := range ch {
				if x.v >= 0 {
					matrix[x.i][x.j] = x.v
				}
				wg.Done()
			}
		}(writeChannel, &writeWaitGroup, &matrix)

		go func(ch chan<- Item, cache *[N][N]float64) {
			for i := 0; i < N; i++ {
				for j := 0; j < N; j++ {
					ch <- Item{i, j, cache[i][j]}
				}
			}
			close(ch)
		}(writeChannel, &cache)

		writeWaitGroup.Wait()
		printMatrix()
	}
}
