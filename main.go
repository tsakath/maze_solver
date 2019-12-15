package main

import (
	"image"
	"fmt"
	_ "image/png"
	"log"
	"os"
)

type Neibor int

const (
	UP Neibor = iota
	RIGHT
	DOWN
	LEFT
)

type Position struct {
	X int // X coordinate
	Y int // Y coordinate
	W int // W weight between neighbours
}

type Node struct {
	Position Position
	Neighbours Neighbours
}

type Neighbours map[Neibor]*Node

func NewNode (x, y int) *Node {
	p := Position{X: x, Y: y, W: 0}
	n := map[Neibor]*Node{}
	return &Node{Position: p, Neighbours: n}
}

func main() {
	reader, err := os.Open("mazes/tiny.png")
	if err != nil {
		log.Fatal(err)
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()

	topRow := [10]*Node{}

	startNode := &Node{}
	endNode := &Node{}

	//var start, end []int
	var previous, current, next, up, down bool
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		leftNode := NewNode(0, y)
		for x := bounds.Min.X; x < bounds.Max.X-1; x++ {
			//fmt.Printf("X: %d, y: %d\n", x, y)
			previous = current
			current = isPath(x, y, m)
			next = isPath(x+1, y, m)
			up = isPath(x, y-1, m)
			down = isPath(x, y+1, m)

			// On wall
			if !current {
				fmt.Print("-")
				continue
			}

			var n *Node

			if previous {
				if next {
					// PATH PATH PATH
					if up || down {
						n = NewNode(x,y)
						leftNode.Neighbours[RIGHT] = n
						n.Neighbours[LEFT] = leftNode
						//fmt.Println("Path above or below", x, y)
					}
				} else {
					// PATH PATH WALL
					n = NewNode(x,y)
					leftNode.Neighbours[RIGHT] = n
					n.Neighbours[LEFT] = leftNode
					//fmt.Println("End of path", x, y)
				}
			} else {
				if next {
					// WALL PATH PATH
					n = NewNode(x, y)
					leftNode = n
					//fmt.Println("Start of path", x, y)
				} else {
					// WALL PATH WALL
					if !up || !down {
						n = NewNode(x,y)
						if y == bounds.Min.Y {
							startNode = n
							// Add startNode to topRow
							topRow[x] = startNode
						}
						if y == bounds.Max.Y-1 {
							endNode = n
						}
						//fmt.Println("Dead end", x, y)
					}
				}
			}

			if n != nil {
				fmt.Printf("O")
				if up {
					// Connect nodes
					n.Neighbours[UP] = topRow[x]
				}
				if down {
					// Add node for the next connection
					topRow[x] = n
				}
			} else {
				fmt.Print("-")
			}
		//	if n != nil && len(n.Neighbours) != 0 {
		//		fmt.Printf("UP%+v\n", n.Neighbours[UP])
		//		fmt.Printf("DOWN%+v\n", n.Neighbours[DOWN])
		//		fmt.Printf("LEFT%+v\n", n.Neighbours[LEFT])
		//		fmt.Printf("RIGHT%+v\n", n.Neighbours[RIGHT])
		//		fmt.Println("")
		//	}
		}
		fmt.Println()
	}
	fmt.Println("startNode", startNode)
	fmt.Println("endNode", endNode)
}

func isPath(x, y int, m image.Image) bool {
	r, _, _, _ := m.At(x, y).RGBA()
	return r > 0
}
