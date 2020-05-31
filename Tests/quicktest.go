package main

import "fmt"

func main() {
	x := 10
	switch {
	case x < 20:
		fmt.Println("< 20")
	case x < 40:
		fmt.Println("< 40")
		
	}
}
