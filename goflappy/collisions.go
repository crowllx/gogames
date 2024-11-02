package main

import (
	"image"
)

type Vertex struct {
	X, Y int
}

type Shape interface{}

type Circle struct {
	cx, cy, r int
}

type ConvexPolygon struct {
}

type BoundingBox struct {
	bb    image.Rectangle
	shape Shape
}
