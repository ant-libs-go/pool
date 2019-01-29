/* ######################################################################
# Author: (zfly1207@126.com)
# Created Time: 2019-01-11 12:30:40
# File Name: pool.go
# Description:
####################################################################### */

package pool

import (
	"container/list"
	"sync"
	"time"
)

// Examples:
//
//	func newPool() *pool.Pool {
//		return &pool.Pool{
//			New: func(key string) interface{} {
//				return NewStruct(key)
//			}
//		}
//	}
//
//  var (
//		Default *pool.Pool = newPool()
//  )
//
//  func main() {
//		cli := Default.Get().Use().(*Struct).P()
//		cli.Close()
//  }
//

var DefaultKey = "default"

type Pool struct {
	New       func(key string) interface{}
	Ping      func(entry Entry) error // 未实现
	MaxActive int                     // 最大连接数，0为不限制，未实现
	MaxIdle   int                     // 最大等待连接数，0为不限制，未实现
	MaxWait   time.Duration           // 最大建立实例时间，0为不限制，未实现
	entries   map[string]*list.List
	lock      sync.RWMutex
}

func New(fn func(key string) interface{}) *Pool {
	o := &Pool{New: fn}
	return o
}

func (this *Pool) Get() *Entry {
	return this.GetByKey(DefaultKey)
}

func (this *Pool) GetByKey(key string) (r *Entry) {
	this.lock.Lock()
	l := this.entries[key]
	if l != nil && l.Len() > 0 {
		ele := l.Back()
		r = ele.Value.(*Entry)
		l.Remove(ele)
	}
	this.lock.Unlock()
	if r == nil {
		r = &Entry{entity: this.New(key), key: key, pool: this}
	}
	return
}

func (this *Pool) Put(entry *Entry) (err error) {
	this.lock.Lock()
	if this.entries == nil {
		this.entries = map[string]*list.List{}
	}
	if _, ok := this.entries[entry.key]; !ok {
		this.entries[entry.key] = list.New()
	}
	this.entries[entry.key].PushFront(entry)
	this.lock.Unlock()
	return
}
