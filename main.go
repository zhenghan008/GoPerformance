package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
)

var (
	i uint64 = 0 // 全局变量i

)

// Consumer 消费者
func Consumer(num int, ch chan uint64) {
	fmt.Println(fmt.Sprintf("%d 开始消费数据", num))
	for n := range ch {
		fmt.Println(fmt.Sprintf("由消费者%d消费的%d出队列", num, n))
	}

}

// Producer 生产者
func Producer(num int, ch chan uint64) {

	fmt.Println(fmt.Sprintf("%d 开始生产数据", num))
	for {
		if i > 1000 {
			break
		}
		ch <- i
		fmt.Println(fmt.Sprintf("由生产者%d生产的%d进入队列", num, i))
		atomic.AddUint64(&i, 1)
	}

}

func main() {
	ch := make(chan uint64, 1000)
	for i := 0; i < 4; i++ {
		go Producer(i, ch)
	}
	for i := 0; i < 3; i++ {
		go Consumer(i, ch)
	}
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("err: ", err)
		}
	}()
	defer close(ch)
	defer fmt.Println("关闭ch")

	//阻塞main函数，ctr + C 退出
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}
