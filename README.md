# Pool

一个简单的对象池

[![License](https://img.shields.io/:license-apache%202-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GoDoc](https://godoc.org/github.com/ant-libs-go/pool?status.png)](http://godoc.org/github.com/ant-libs-go/pool)
[![Go Report Card](https://goreportcard.com/badge/github.com/ant-libs-go/pool)](https://goreportcard.com/report/github.com/ant-libs-go/pool)

# 特性

* 支持最大空闲、最大缓存对象数限制
* 支持无可用对象时阻塞等待
* 支持按对象使用次数或时间进行销毁
* 支持对象数据使用前reset对象内数据

# 快速开始

```golang
type SearchContext struct {
	No    int32
	Param int32
}

func (this *SearchContext) Reset() {
	this.Param = 0
}

p := &Pool{
    New: func() Obj {
        return &SearchContext{No: i}
    },
    Test: func(o Obj) bool {
        if o.(*SearchContext).No == 2 {
            return false
        }
        return true
    },
    MaxActive: 2,
    Wait:      true,
}

// 获取对象
entry := p.Get()
// 使用对象
fmt.Println("---", entry.Obj().Param)
// 将对象重置到缓存池
entry.Close()
```
