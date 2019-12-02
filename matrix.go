package main

import (
	"math"
)

// [rows][cols]float64
type Matrix [3][3]float64

func (m Matrix) Apply(v Point) Point {
	return Point{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2],
	}
}

func (m Matrix) Mul(w Matrix) Matrix {
	return Matrix{
		{
			m[0][0]*w[0][0] + m[0][1]*w[1][0] + m[0][2]*w[2][0],
			m[0][0]*w[0][1] + m[0][1]*w[1][1] + m[0][2]*w[2][1],
			m[0][0]*w[0][2] + m[0][1]*w[1][2] + m[0][2]*w[2][2],
		},
		{
			m[1][0]*w[0][0] + m[1][1]*w[1][0] + m[1][2]*w[2][0],
			m[1][0]*w[0][1] + m[1][1]*w[1][1] + m[1][2]*w[2][1],
			m[1][0]*w[0][2] + m[1][1]*w[1][2] + m[1][2]*w[2][2],
		},
		{
			m[2][0]*w[2][0] + m[2][1]*w[1][0] + m[2][2]*w[2][0],
			m[2][0]*w[2][1] + m[2][1]*w[1][1] + m[2][2]*w[2][1],
			m[2][0]*w[2][2] + m[2][1]*w[1][2] + m[2][2]*w[2][2],
		},
	}
}

func NewRotationX(theta float64) Matrix {
	c := math.Cos(theta)
	s := math.Sin(theta)

	return Matrix{
		{1, 0, 0},
		{0, c, -s},
		{0, s, c},
	}
}

func NewRotationY(theta float64) Matrix {
	c := math.Cos(theta)
	s := math.Sin(theta)

	return Matrix{
		{c, 0, s},
		{0, 1, 0},
		{-s, 0, c},
	}
}

func NewRotationZ(theta float64) Matrix {
	c := math.Cos(theta)
	s := math.Sin(theta)

	return Matrix{
		{c, -s, 0},
		{s, c, 0},
		{0, 0, 1},
	}
}
