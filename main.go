package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var ProducerMutex sync.Mutex
var i = 0

func Consumer(num int, ch chan int) {
	fmt.Println(fmt.Sprintf("%d 开始消费数据", num))
	for {

		fmt.Println(fmt.Sprintf("由消费者%d消费的%d出队列", num, <-ch))
		time.Sleep(10 * time.Microsecond)

	}

}

func Producer(num int, ch chan int) {

	fmt.Println(fmt.Sprintf("%d 开始生产数据", num))
	for {
		if i > 1000 {
			break
		}
		ProducerMutex.Lock()
		ch <- i
		fmt.Println(fmt.Sprintf("由生产者%d生产的%d进入队列", num, i))
		i += 1
		ProducerMutex.Unlock()
		time.Sleep(10 * time.Microsecond)
	}
}

func main() {
	ch := make(chan int, 3000)
	for i := 1; i <= 3; i++ {
		go Producer(i, ch)
		go Consumer(i, ch)

	}

	defer close(ch)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}
