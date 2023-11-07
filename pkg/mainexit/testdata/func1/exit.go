package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("main()")
	HiddenExit(1)
}

func HiddenExit(code int) {
	fmt.Println("HiddenExit()")

	os.Exit(code)
}
