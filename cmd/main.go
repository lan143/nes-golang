package main

import (
	"context"
	"fmt"
	"main/src"
	_ "net/http/pprof"
	"os"
)

func main() {
	container := src.BuildContainer()
	err := container.Invoke(func(nes *src.Nes) {
		if len(os.Args) == 1 || len(os.Args[1]) == 0 {
			fmt.Println("Nintendo Entertainment System (NES) emulator written in Golang\n")
			fmt.Println("Usage:\n")
			fmt.Printf("        %s <rom_name>\n\n", os.Args[0])
			fmt.Println("Where <rom_name> is a path to NES ROM\n")
			return
		}

		err := nes.Init(os.Args[1])
		if err != nil {
			panic(err)
		}

		nes.Run(context.Background())
	})

	if err != nil {
		panic(err)
	}
}
