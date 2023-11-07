package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("main")
	os.Exit(0) // want `os.Exit called in main func in main package`
}
