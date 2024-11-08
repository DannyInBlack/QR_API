package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
)

// This program implements the Cantor Pairing function and it's inverse
// Danial Sabry 211010447

// Makes sure that command line arguments are in {1, 2}
func formatArgs(args []string) (int, error) {
	if len(args) == 0 {
		return -1, errors.New("no arguments given")
	}

	if len(args) == 1 {
		if ty, err := strconv.Atoi(args[0]); err == nil {
			if ty != 1 && ty != 2 {
				return -1, errors.New("invalid arguments")
			}
			return ty, err
		} else {
			return -1, err
		}
	}

	return -1, errors.New("too many arguments given")
}

// Cantor
func foo(x int, y int) int {
	return (x+y+1)*(x+y)/2 + y
}

// Inverse Cantor
func fooInverse(z int) (int, int) {
	w := int((math.Sqrt(float64(8*z+1)) - 1) / 2)
	t := (w*w + w) / 2
	x, y := w-z+t, z-t

	return x, y
}

// Wrapper that handles incorrect formatting and input for Cantor
func cantor() {
	x, y := -1, -1

	for x <= 0 && y <= 0 {
		fmt.Println("Please enter two positive integers: X and Y")

		if _, err := fmt.Scanf("%d %d", &x, &y); err != nil {
			fmt.Println("Error: ", err)
		}

	}

	fmt.Println(foo(x, y))
}

// Wrapper that handles incorrect formatting and input for Inverse Cantor
func inverseCantor() {
	z := -1

	for z <= 0 {
		fmt.Println("Please enter a positive integer: Z")

		if _, err := fmt.Scanf("%d", &z); err != nil {
			fmt.Println("Error: ", err)
		}

	}

	fmt.Println(fooInverse(z))
}

func main() {

	args, err := formatArgs(os.Args[1:])

	if err != nil {
		// Exit with error from input args
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	if args == 1 {
		cantor()
	} else {
		inverseCantor()
	}

}
