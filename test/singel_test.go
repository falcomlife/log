package test

import (
	"fmt"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {
	fmt.Println("testing...")
	ch := make(chan int)
	go func() {
		ch <- 1
		ch <- 2
		ch <- 3
		ch <- 4
		ch <- 5
		ch <- 6
	}()
	time.Sleep(10 * time.Second)
	fmt.Println(ch)
	go func() {
		select {
		case num := <-ch:
			fmt.Println(num)
		}
	}()
}
