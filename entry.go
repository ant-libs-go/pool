/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2019-01-11 13:09:36
# File Name: entry.go
# Description:
####################################################################### */

package pool

import (
	"sync"
)

type Entry struct {
	key    string
	entity interface{}
	pool   *Pool
	use    bool
	count  int
	lock   sync.RWMutex
}

func (this *Entry) Use() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.count++
	this.use = true
	return this.entity
}

func (this *Entry) Close() {
	this.use = false
	this.pool.Put(this)
}
