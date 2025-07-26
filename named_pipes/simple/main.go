package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	coffeeShopPipePath := "/tmp/coffee_shop_pipe"
	if !namedPipeEixsts(coffeeShopPipePath) {
		println("no named pipe found at: ", coffeeShopPipePath)
		if err := unix.Mkfifo(coffeeShopPipePath, 0666); err != nil {
			fmt.Printf("errror while creating pipe: %v\n", err)
			return
		}
		println("file created")
	}
}

func namedPipeEixsts(pipePath string) bool {
	info, err := os.Stat(pipePath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeNamedPipe != 0
}
