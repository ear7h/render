package main

import "testing"

func TestMatrix(t *testing.T) {
	type tcase struct {
		a, b, c Matrix
		f       func(a, b Matrix) Matrix
	}

	fn := func(tc tcase) func(t *testing.T) {
		return func(t *testing.T) {
			c := tc.f(tc.a, tc.b)
			if c != tc.c {
				t.Fatalf("incorect restult:\n%v\nexpected:\n%v", c, tc.c)
			}
		}
	}

	tcases := map[string]tcase{
		"mul": {
			a: Matrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			b: Matrix{
				{1, 0, 0},
				{1, 0, 0},
				{1, 0, 0},
			},
			c: Matrix{
				{1, 0, 0},
				{1, 0, 0},
				{1, 0, 0},
			},
			f: Matrix.Mul,
		},
	}

	for k, v := range tcases {
		t.Run(k, fn(v))
	}
}
