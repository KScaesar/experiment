package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func A(wg *sync.WaitGroup) chan string {
	c := make(chan string)

	go func() {
		rand.Seed(time.Now().Unix())
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // 模擬某種 io 讀取

		defer wg.Done()
		c <- "A_function"
	}()
	return c
}

func B(wg *sync.WaitGroup) chan string {
	c := make(chan string)

	go func() {
		rand.Seed(time.Now().Unix())
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // 模擬某種 io 讀取

		defer wg.Done()
		c <- "B_function"
	}()
	return c
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(2)
	aChan := A(&wg)
	bChan := B(&wg)

	aResult := make(chan string)
	bResult := make(chan string)
	go func() {
		for v := range aChan {
			aResult <- v
		}
	}()

	go func() {
		for v := range bChan {
			bResult <- v
		}
	}()

	wg.Wait()
	close(aChan)
	close(bChan)

	fmt.Printf("a=%v b=%v", <-aResult, <-bResult)
	// 若沒有註解 sleep
	// 可能出現三種情況
	// 1. a= b=B_function
	// 2. a=A_function b=
	// 3. a=A_function b=B_function
	//
	// 若註解 sleep
	// 只會出現
	// a=A_function b=B_function
	// https://play.golang.org/p/xLjrT8oTvdi
	//
	// 不使用 range, 改成普通的 chan 接收
	// https://play.golang.org/p/aiZP4ifkPoR
	//
	// 改成使用 resultChan 才能保證接收到數值
	// https://play.golang.org/p/lggis2dGM9T
	//
	// 相關議題
	// 記憶體同步化
	// happens before
}
