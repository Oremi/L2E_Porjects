package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// SolveQuadratic returns the roots of ax^2 + bx + c = 0.

func SolveQuadratic(a, b, c float64) (x1, x2 float64, isComplex bool, err error) {

	if a == 0 && b == 0 {
		if c == 0 {
			return 0, 0, false, fmt.Errorf("infinite solutions (0=0)")
		}
		return 0, 0, false, fmt.Errorf("no solution (constant != 0)")
	}
	if a == 0 {
		// Linear equation: bx + c = 0 => x = -c/b
		return -c / b, 0, false, nil
	}

	discriminant := b*b - 4*a*c

	if discriminant >= 0 {
		sqrtDiscriminant := math.Sqrt(discriminant)
		x1 = (-b - sqrtDiscriminant) / (2 * a)
		x2 = (-b + sqrtDiscriminant) / (2 * a)
		return x1, x2, false, nil
	}
	// the Complex roots of the equation
	realPart := -b / (2 * a)
	imagPart := math.Sqrt(-discriminant) / (2 * a)
	return realPart, imagPart, true, nil
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run quadratic.go a b c")
		fmt.Println("Example: go run quadratic.go 1 -3 2")
		return
	}
	a, err1 := strconv.ParseFloat(os.Args[1], 64)
	b, err2 := strconv.ParseFloat(os.Args[2], 64)
	c, err3 := strconv.ParseFloat(os.Args[3], 64)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("Error: all arguments must be numbers")
		return
	}

	x1, x2, isComplex, err := SolveQuadratic(a, b, c)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if isComplex {
		fmt.Printf("Roots are complex: %.3f ± %.3fi\n", x1, x2)
	} else if x1 == x2 {
		fmt.Printf("Double root: %.3f\n", x1)
	} else {
		fmt.Printf("Two real roots: %.3f and %.3f\n", x1, x2)
	}
}
