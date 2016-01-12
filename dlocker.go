package DLocker

import (
	"github.com/nladuo/DLocker/model"
	"github.com/samuel/go-zookeeper/zk"
	"log"
)

type Dlocker struct {
	lockerPath string
	prefix     string
	basePath   string
}

func NewLocker(path string, prefix string) *Dlocker {

	var locker Dlocker
	locker.basePath = path
	locker.prefix = prefix
	isExsit, _, err := getZkConn().Exists(path)
	if err != nil {
		panic(err.Error())
	}
	if !isExsit {
		log.Println("create the znode:" + path)
		getZkConn().Create(path, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
	} else {
		log.Println("the znode " + path + " existed")
	}
	return &locker
}

func (this *Dlocker) Lock() (isSuccess bool) {
	isSuccess = false
	defer func() {
		e := recover()
		if e == zk.ErrConnectionClosed {
			//try reconnect the zk server
			log.Println("connection closed, reconnect to the zk server")
			reConnectZk()
		}
	}()
	//create znode
	path := this.basePath + "/" + this.prefix
	var err error
	this.lockerPath, err = getZkConn().Create(path, []byte(""), zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
	//get children and check is the created locker is the minimum znode
	chidren, _, err := getZkConn().Children(this.basePath)
	if err != nil {
		panic(err)
	}
	minIndex := model.GetMinIndex(chidren, this.prefix)
	minLockerPath := this.basePath + "/" + chidren[minIndex]
	// if the created znode is not the minimum znode,
	// listen for the pre-znode delete notification
	if minLockerPath != this.lockerPath {
		lastNodeName := model.GetLastNodeName(this.lockerPath,
			this.basePath, this.prefix)
		isExist, _, watch, err := getZkConn().ExistsW(this.basePath + "/" + lastNodeName)
		if err != nil {
			panic(err)
		}
		if isExist {
			event := <-watch
			if event.Type == zk.EventNodeDeleted {
				isSuccess = true
			} else {
				isSuccess = false
			}
		}

	} else { // if the created node is the minimum znode, getLock success
		isSuccess = true
	}
	return
}

func (this *Dlocker) Unlock() (isSuccess bool) {
	isSuccess = false
	defer func() {
		e := recover()
		if e == zk.ErrConnectionClosed {
			//try reconnect the zk server
			log.Println("connection closed, reconnect to the zk server")
			reConnectZk()
		}
	}()
	err := getZkConn().Delete(this.lockerPath, 0)
	if err != nil {
		panic(err)
	}
	isSuccess = true
	return
}
