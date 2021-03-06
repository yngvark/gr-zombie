package main

import (
"context"
"fmt"
"time"
)

func main2() {
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx2 context.Context) {
		fmt.Println("THREAD STARTING")

		ticker := time.NewTicker(time.Second * 1) //nolint:gomnd
		defer ticker.Stop()

		a := 0
		for true {
			select {
			case <-ctx2.Done():
				fmt.Println("THREAD DONE!")
				return
			case <-ticker.C:
				fmt.Printf("THREAD-i = %d\n", a)
				a++
			}
		}
	}(ctx)

	fmt.Println("MAIN FUNCTION COUNTING")

	i := 0
	for i < 2 {
		fmt.Printf("MAIN-i = %d\n", i)
		i++
		time.Sleep(1 * time.Second)
	}

	cancel()

	fmt.Println("MAIN WAITING FOR STUFF")
	time.Sleep(2 * time.Second)
}

func a() {
	a := make(chan int)
	a <- 3

	//b := make(chan<- struct{})
	//b <- struct{}{}
	//
	//c := <- b
	//fmt.Println(c)
}

type arne struct {}

