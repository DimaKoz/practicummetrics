package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("main()")
	t := Finisher{}
	t.Done(1)
}

type Finisher struct {
}

func (e Finisher) Done(code int) {
	fmt.Println("Done()")
	os.Exit(code)
}
