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
	"sync/atomic"
	"time"
)

type Obj interface {
	Reset()
}

type Entry struct {
	obj   Obj
	pool  *Pool
	usage int32
	ctime time.Time
}

func (this *Entry) isValid(maxUsage int32, maxLifetime time.Duration) bool {
	if maxUsage != 0 && this.usage >= maxUsage {
		return false
	}
	if maxLifetime != 0 && time.Now().Sub(this.ctime) >= maxLifetime {
		return false
	}
	return true
}

func (this *Entry) use() *Entry {
	atomic.AddInt32(&this.usage, 1)
	return this
}

func (this *Entry) Obj() Obj {
	return this.obj
}

func (this *Entry) Close() {
	this.pool.put(this)
}

type Pool struct {
	lock    sync.RWMutex
	active  int32 // 当前实例化的对象数
	entries *list.List

	New         func() Obj
	Test        func(Obj) bool
	MaxIdle     int32         // 最大空闲对象数，0为不限制。超出限制的对象Put时将丢弃
	MaxActive   int32         // 最大缓存对象数，0为不限制。超出限制将无实例返回
	MaxUsage    int32         // 对象最大使用次数，0为不限制。超出限制将丢弃
	MaxLifetime time.Duration // 对象最大生存时间，0为不限制。超出限制将丢弃
	Wait        bool          // 当无可用对象返回时，Get是否等待
}

func New(fn func() Obj) *Pool {
	o := &Pool{New: fn}
	return o
}

func (this *Pool) Active() int32 {
	return this.active
}

func (this *Pool) Idle() int32 {
	return int32(this.entries.Len())
}

func (this *Pool) Get() *Entry {
	this.lock.Lock()
	if this.entries == nil {
		this.entries = list.New()
	}
	this.lock.Unlock()

	var entry *Entry

RETRY:
	for this.Idle() > 0 {
		entry = this.entries.Remove(this.entries.Back()).(*Entry)

		if entry.isValid(this.MaxUsage, this.MaxLifetime) == false {
			entry = nil
			atomic.AddInt32(&this.active, -1)
			continue
		}
		if this.Test != nil && this.Test(entry.Obj()) == false {
			entry = nil
			atomic.AddInt32(&this.active, -1)
			continue
		}
		break
	}
	if entry == nil && (this.MaxActive == 0 || this.MaxActive > this.Active()) {
		atomic.AddInt32(&this.active, 1)
		entry = &Entry{obj: this.New(), ctime: time.Now(), pool: this}
	}
	if entry == nil && this.Wait == true {
		time.Sleep(time.Microsecond * 10)
		goto RETRY
	}

	if entry == nil {
		return nil
	}
	return entry.use()
}

func (this *Pool) put(entry *Entry) (err error) {
	this.lock.Lock()
	if this.entries == nil {
		this.entries = list.New()
	}
	this.lock.Unlock()

	if this.MaxIdle != 0 && this.MaxIdle <= this.Idle() {
		atomic.AddInt32(&this.active, -1)
		return
	}
	if entry.isValid(this.MaxUsage, this.MaxLifetime) == false {
		atomic.AddInt32(&this.active, -1)
		return
	}

	this.entries.PushFront(entry)
	return
}
