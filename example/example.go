package main

import (
	"fmt"
	"github.com/nladuo/DLocker"
	"log"
	"time"
)

var (
	hosts         []string      = []string{"127.0.0.1:2181"} // the zookeeper hosts
	basePath      string        = "/locker"                  //the application znode path
	lockerTimeout time.Duration = 1 * time.Minute            // the maximum time for a locker waiting
	zkTimeOut     time.Duration = 20 * time.Second           // the zk connection timeout
)

func main() {
	end := make(chan byte)
	err := DLocker.EstablishZkConn(hosts, zkTimeOut)
	defer DLocker.CloseZkConn()
	if err != nil {
		panic(err)
	}
	//concurrent test lock and unlock
	for i := 0; i < 100; i++ {
		go run(i)
	}
	<-end
}

func run(i int) {
	locker := DLocker.NewLocker(basePath, lockerTimeout)
	for {
		locker.Lock() // like mutex.Lock()
		fmt.Println("gorountine", i, ": get lock")
		//do something of which time not excceed lockerTimeout
		fmt.Println("gorountine", i, ": unlock")
		if !locker.Unlock() { // like mutex.Unlock()
			log.Println("gorountine", i, ": unlock failed")
		}
	}
}
