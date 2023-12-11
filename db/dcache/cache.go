/*
File Name:  dcache.go
Description:
Author:      Chenghu
Date:       2023/8/23 22:21
Change Activity:
*/
package dcache

import (
	"github.com/fanjindong/go-cache"
	"time"
)

var GoCache BCache

func init() {
	initCache()
}

func initCache() {
	////创建的时候可以设置  过期回调函数
	//f := func(k string, v interface{}) error {
	//	fmt.Println("ExpiredCallback", k, v)
	//	return nil
	//}
	//GoCache = dcache.NewMemCache(dcache.WithExpiredCallback(f))

	ch := cache.NewMemCache()
	GoCache = BCache{ch}
	//c.Set("a", 1)
	//c.Set("b", 1, dcache.WithEx(1*time.Second))
	//c.Get("a") // 1, true
	//c.Get("b") // nil, false
}

type BCache struct {
	cache.ICache
}

// 加锁, 等待锁
// exp ms

func (bc BCache) WaitingLock(name string, exp int) bool {
	for {
		lastExp, ok := bc.Ttl(name)
		if ok {
			time.Sleep(lastExp)
			continue
		}
		break
	}
	if exp == 0 {
		exp = 30
	}
	bc.Set(name, 1, cache.WithEx(time.Duration(exp)*time.Millisecond))
	return true
}

func (bc BCache) Unlock(name string) {
	bc.Del(name)
}
