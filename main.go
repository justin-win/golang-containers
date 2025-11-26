package main

import (
	"fmt"
	"os"

	"gopod/src/engine"
)

func main() {
	if argc := len(os.Args); argc < 2 {
		fmt.Println("./gopod <run> <commands>")
		os.Exit(1)
	}
	switch os.Args[1] {
		case "run":
			engine.SetUpChild()
		case "child":
			if uid := os.Getuid(); uid != 0 {
				fmt.Println("./gopod <run> <commands>")
				os.Exit(-1)
			}
			engine.CreateContainer()
		default:
			fmt.Println("./gopod <run> <commands>")
			os.Exit(1)
	}
}

