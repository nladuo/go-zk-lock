package main

import (
	"fmt"
	"github.com/nladuo/DLocker"
	"log"
	"time"
)

var (
	hosts         []string      = []string{"127.0.0.1:2181"}
	basePath      string        = "/locker"
	prefix        string        = "lock-"
	lockerTimeout time.Duration = 8 * time.Second
	zkTimeOut     time.Duration = 20 * time.Second
)

func run(i int) {
	locker := DLocker.NewLocker(basePath, prefix, lockerTimeout)
	for {
		for !locker.Lock() {
		}
		fmt.Println("gorountine ", i, " get lock")
		time.Sleep(time.Millisecond * 10)
		fmt.Println("gorountine ", i, " unlock")
		if !locker.Unlock() {
			log.Println("gorountine ", i, "unlock failed")
		}

	}

}

func main() {
	ch := make(chan byte)
	err := DLocker.EstablishZkConn(hosts, zkTimeOut)
	defer DLocker.CloseZkConn()
	if err != nil {
		panic(err)
	}
	fmt.Println("hello")
	for i := 0; i < 15; i++ {
		go run(i)
	}

	<-ch

}
