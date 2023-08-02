package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

func main() {
	etl1()
	etl2()
	etl3()
}

func etl1() {
	var workIdAll []int
	for i := 0; i < 9; i++ {
		workIdAll = append(workIdAll, i)
	}

	fanOutAll := lop.Map(lo.Chunk(workIdAll, 2), func(workIdChunk []int, _ int) <-chan []string {
		mq := make(chan []string, 1) // 和 etl2 的差異
		transform := func(workId int, _ int) string {
			return strconv.Itoa(workId)
		}
		fmt.Println(workIdChunk)
		mq <- lo.Map(workIdChunk, transform)
		close(mq)
		return mq
	})

	chunkResult := lo.ChannelToSlice(lo.FanIn(0, fanOutAll...))
	result := lo.Flatten(chunkResult)
	fmt.Println("etl1:", result)
}

func etl2() {
	var workIdAll []int
	for i := 0; i < 9; i++ {
		workIdAll = append(workIdAll, i)
	}

	fanOutAll := lo.Map(lo.Chunk(workIdAll, 2), func(workIdChunk []int, _ int) <-chan []string {
		mq := make(chan []string) // 和 etl1 的差異
		transform := func(workId int, _ int) string {
			return strconv.Itoa(workId)
		}
		go func() {
			fmt.Println(workIdChunk)
			mq <- lo.Map(workIdChunk, transform)
			close(mq)
		}()
		return mq
	})

	chunkResult := lo.ChannelToSlice(lo.FanIn(0, fanOutAll...))
	result := lo.Flatten(chunkResult)
	fmt.Println("etl2:", result)
}

func etl3() {
	var workIdAll []int
	for i := 0; i < 9; i++ {
		workIdAll = append(workIdAll, i)
	}

	fanOutAll := lo.Map(lo.Chunk(workIdAll, 2), func(workIdChunk []int, _ int) <-chan string {
		mq := make(chan string, len(workIdChunk))
		timeout, cancelFunc := context.WithTimeout(context.Background(), time.Minute)
		go func() {
			defer func() {
				close(mq)
				cancelFunc()
			}()

			fmt.Println(workIdChunk)
			lo.ForEach(workIdChunk, func(workId int, index int) {
				select {
				case mq <- strconv.Itoa(workId):
				case <-timeout.Done():
					return
				}
			})
		}()
		return mq
	})

	result := lo.ChannelToSlice(lo.FanIn(0, fanOutAll...))
	fmt.Println("etl3:", result)

}

type User struct {
	address []string
}

var example1 Tuple2[User, error]
var example2 Tuple2[*User, error]

type Tuple2[A any, B any] struct {
	A A
	B B
}
