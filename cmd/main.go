package main

import (
	"context"
	"main/src"
)

func main() {
	container := src.BuildContainer()
	err := container.Invoke(func(nes *src.Nes) {
		err := nes.Init()
		if err != nil {
			panic(err)
		}

		nes.Run(context.Background())
	})

	if err != nil {
		panic(err)
	}
}
