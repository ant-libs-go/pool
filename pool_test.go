/* ######################################################################
# Author: (zhengfei@dianzhong.com)
# Created Time: 2021-10-20 14:30:47
# File Name: pool_test.go
# Description:
####################################################################### */

package pool

import (
	"fmt"
	"testing"
)

type SearchContext struct {
	No    int32
	Param int32
}

func (this *SearchContext) Reset() {
	//this.Param = 0
}

var i int32

func TestBasic(t *testing.T) {
	p := &Pool{
		New: func() Obj {
			i = i + 1
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

	entry1 := p.Get()
	fmt.Println("---", entry1.Obj())
	entry2 := p.Get()
	fmt.Println("---", entry2.Obj())
	entry3 := p.Get()
	fmt.Println("---", entry3.Obj())
	entry4 := p.Get()
	fmt.Println("---", entry4.Obj())
	entry1.Close()
	entry2.Close()
	entry3.Close()
	entry4.Close()

	entry101 := p.Get()
	fmt.Println("---", entry101.Obj())
	entry101.Close()

	entry102 := p.Get()
	fmt.Println("---", entry102.Obj())
	entry102.Close()

	entry103 := p.Get()
	fmt.Println("---", entry103.Obj())
	entry103.Close()

	entry104 := p.Get()
	fmt.Println("---", entry104.Obj())
	entry104.Close()

	// t.Fatal("not implemented")
}

// vim: set noexpandtab ts=4 sts=4 sw=4 :
