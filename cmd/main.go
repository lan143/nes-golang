package main

import (
	"context"
	"main/src"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		http.ListenAndServe(":8081", nil)
	}()

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
