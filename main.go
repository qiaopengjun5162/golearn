package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// 5个协程并发，发送随机数量的随机数到一个channel，另外一个协程从channel读取数据，不能用time.sleep，求和并打印结果，要求空间复杂度为O(1)。

func producer(id int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := rand.Intn(10) + 1 // 随机1-10个数
	for i := 0; i < count; i++ {
		num := rand.Intn(100) // 随机0-99的数
		ch <- num
		fmt.Printf("Producer %d sent: %d\n", id, num)
	}
}

func consumer(ch <-chan int, done chan<- bool) {
	sum := 0
	for num := range ch {
		sum += num
		fmt.Printf("Consumer received: %d, current sum: %d\n", num, sum)
	}
	fmt.Println("Final sum:", sum)
	done <- true
}

func main() {
	const producerCount = 5
	ch := make(chan int)
	var wg sync.WaitGroup
	done := make(chan bool)

	// 启动消费者
	go consumer(ch, done)

	// 启动生产者
	wg.Add(producerCount)
	for i := 0; i < producerCount; i++ {
		go producer(i, ch, &wg)
	}

	// 等待所有生产者完成
	go func() {
		wg.Wait()
		close(ch) // 关闭channel通知消费者
	}()

	// 等待消费者完成
	<-done
}
