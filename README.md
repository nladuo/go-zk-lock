# go-zk-lock
[![Build Status](https://travis-ci.org/nladuo/go-zk-lock.svg)](https://travis-ci.org/nladuo/go-zk-lock)

a distributed locker based on zookeeper and implemented in golang.

# Installation
```go
go get github.com/samuel/go-zookeeper/zk
go get github.com/nladuo/go-zk-lock
```

# Usage
#### Configure the Zookeeper
You can check out the zookeeper configuration <a href="http://zookeeper.apache.org/doc/r3.4.6/zookeeperStarted.html">here</a>.
#### Set Your Configuration
``` go
var (
        hosts         []string      = []string{"127.0.0.1:2181"} // the zookeeper hosts
        basePath      string        = "/locker"                  //the application znode path
        lockerTimeout time.Duration = 1 * time.Minute            // the maximum time for a locker waiting
        zkTimeOut     time.Duration = 20 * time.Second           // the zk connection timeout
)
```
#### Establish Zookeeper Connection
``` go
err := DLocker.EstablishZkConn(hosts, zkTimeOut)
defer DLocker.CloseZkConn()
if err != nil {
        panic(err)
}
```
#### Create Distributed Locker
``` go
locker := DLocker.NewLocker(basePath, lockerTimeout)
```
#### Lock And Unlock
``` go
locker.Lock() // like mutex.Lock()
//do something of which time not excceed lockerTimeout
if !locker.Unlock() { // like mutex.Unlock(), return false when zookeeper connection error or locker timeout
        log.Println("Sorry, unlock failed")
}
```
