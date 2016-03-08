package DLocker

import (
	"github.com/nladuo/DLocker/modules"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"sync"
	"time"
)

type Dlocker struct {
	lockerPath string
	prefix     string
	basePath   string
	timeout    time.Duration
	innerLock  *sync.Mutex
}

func NewLocker(path string, timeout time.Duration) *Dlocker {

	var locker Dlocker
	locker.basePath = path
	locker.prefix = "lock-" //the prefix of a znode, everything is okay
	locker.timeout = timeout
	locker.innerLock = &sync.Mutex{}
	isExsit, _, err := getZkConn().Exists(path)
	locker.checkErr(err)
	if !isExsit {
		log.Println("create the znode:" + path)
		getZkConn().Create(path, []byte(""), int32(0), zk.WorldACL(zk.PermAll))
	} else {
		log.Println("the znode " + path + " existed")
	}
	return &locker
}

func (this *Dlocker) createZnodePath() (string, error) {
	path := this.basePath + "/" + this.prefix
	return getZkConn().Create(path, []byte(""), zk.FlagSequence|zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
}

//get the path of minimum serial number znode from sequential children
func (this *Dlocker) getMinZnodePath() (string, error) {
	children, err := this.getPathChildren()
	if err != nil {
		return "", err
	}
	minSNum := modules.GetMinSerialNumber(children, this.prefix)
	minZnodePath := this.basePath + "/" + children[minSNum]
	return minZnodePath, nil
}

//get the children of basePath znode
func (this *Dlocker) getPathChildren() ([]string, error) {
	children, _, err := getZkConn().Children(this.basePath)
	return children, err
}

//get the last znode of created znode
func (this *Dlocker) getLastZnodePath() string {
	return modules.GetLastNodeName(this.lockerPath,
		this.basePath, this.prefix)
}

//just list mutex.Lock()
func (this *Dlocker) Lock() {
	for !this.lock() {
	}
}

//just list mutex.Unlock()
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
	if err == zk.ErrNoNode {
		isSuccess = false
		return
	} else {
		this.checkErr(err)
	}
	isSuccess = true
	return
}

func (this *Dlocker) lock() (isSuccess bool) {
	isSuccess = false
	defer func() {
		e := recover()
		if e == zk.ErrConnectionClosed {
			//try reconnect the zk server
			log.Println("connection closed, reconnect to the zk server")
			reConnectZk()
		}
	}()
	this.innerLock.Lock()
	defer this.innerLock.Unlock()
	//create a znode for the locker path
	var err error
	this.lockerPath, err = this.createZnodePath()
	this.checkErr(err)

	//get the znode which get the lock
	minZnodePath, err := this.getMinZnodePath()
	this.checkErr(err)

	if minZnodePath == this.lockerPath {
		// if the created node is the minimum znode, getLock success
		isSuccess = true
	} else {
		// if the created znode is not the minimum znode,
		// listen for the pre-znode delete notification
		lastNodeName := this.getLastZnodePath()
		watchPath := this.basePath + "/" + lastNodeName
		isExist, _, watch, err := getZkConn().ExistsW(watchPath)
		this.checkErr(err)
		if isExist {
			select {
			//get lastNode been deleted event
			case event := <-watch:
				if event.Type == zk.EventNodeDeleted {
					//check out the lockerPath existence
					isExist, _, err = getZkConn().Exists(this.lockerPath)
					this.checkErr(err)
					if isExist {
						//checkout the minZnodePath is equal to the lockerPath
						minZnodePath, err := this.getMinZnodePath()
						this.checkErr(err)
						if minZnodePath == this.lockerPath {
							isSuccess = true
						}
					}
				}
			//time out
			case <-time.After(this.timeout):
				// if timeout, delete all the node less than created locker node
				children, err := this.getPathChildren()
				this.checkErr(err)
				deletePathList := modules.GetPathListSerialNumberLessThanLocker(children, this.basePath, this.prefix, this.lockerPath)
				for i := 0; i < len(deletePathList); i++ {
					deletePath := deletePathList[i]
					log.Println("timeout,delete", deletePath, "")
					getZkConn().Delete(deletePath, 0)
				}
			}
		}
	}
	return
}

func (this *Dlocker) checkErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
