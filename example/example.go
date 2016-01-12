package main

import (
	"fmt"
	"github.com/nladuo/DLocker"
	"time"
)

var (
	hosts    []string = []string{"127.0.0.1:2181"}
	basePath string   = "/locker"
	prefix   string   = "lock-"
)

func run(i int) {
	locker := DLocker.NewLocker(basePath, prefix)
	for {
		locker.Lock()
		fmt.Println("gorountine ", i, " get lock")
		time.Sleep(time.Millisecond * 100)
		fmt.Println("gorountine ", i, " unlock")
		locker.Unlock()
	}

}

func main() {
	ch := make(chan byte)
	err := DLocker.EstablishZkConn(hosts)
	if err != nil {
		panic(err)
	}
	fmt.Println("hello")
	for i := 0; i < 15; i++ {
		go run(i)
	}

	<-ch

}
