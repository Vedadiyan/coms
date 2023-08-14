package main

import (
	"fmt"
	"os"

	flaggy "github.com/vedadiyan/flaggy/pkg"
)

func main() {
	err := flaggy.Parse(&Options{}, os.Args)
	if err != nil {
		panic(err)
	}
	fmt.Scanln()
}
