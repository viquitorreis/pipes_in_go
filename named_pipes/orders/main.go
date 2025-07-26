package main

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	coffeeShopPipePath := "/tmp/coffee_shop_pipe"
	if !namedPipeExists(coffeeShopPipePath) {
		if err := unix.Mkfifo(coffeeShopPipePath, 0666); err != nil {
			fmt.Printf("errror while creating pipe: %v\n", err)
			return
		}
	}

	coffeeOrders, err := os.OpenFile(coffeeShopPipePath, os.O_RDWR, os.ModeNamedPipe)
	if err != nil {
		fmt.Printf("error while opening named pipe: %+v", err)
		return
	}
	defer coffeeOrders.Close()

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go produceOrders(wg, coffeeOrders, []string{"Espresso", "Arabica", "Latte", "Cappuccino"})
	go consumeOrders(wg, coffeeOrders)
	wg.Wait()
}

func namedPipeExists(pipePath string) bool {
	info, err := os.Stat(pipePath)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeNamedPipe != 0
}

func produceOrders(wg *sync.WaitGroup, coffeOrders *os.File, orders []string) {
	defer wg.Done()
	for _, order := range orders {
		_, err := coffeOrders.WriteString(order + "\n")
		if err != nil {
			fmt.Printf("err while writing to named pipe: %+v", err)
		}
	}
	coffeOrders.WriteString("finish\n")
}

func consumeOrders(wg *sync.WaitGroup, coffeeOrders *os.File) {
	defer wg.Done()
	scanner := bufio.NewScanner(coffeeOrders)
	for scanner.Scan() {
		order := scanner.Text()
		if order == "finish" {
			println("finishing coffee orders.")
			break
		}
		time.Sleep(time.Second * 2)
		fmt.Println("Processing order: ", order)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading from named pipe: ", err.Error())
	}
}
