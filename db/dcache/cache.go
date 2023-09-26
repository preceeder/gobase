/*
File Name:  dcache.py
Description:
Author:      Chenghu
Date:       2023/8/23 22:21
Change Activity:
*/
package dcache

import (
	"github.com/fanjindong/go-cache"
)

var GoCache cache.ICache

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
	GoCache = ch
	//c.Set("a", 1)
	//c.Set("b", 1, dcache.WithEx(1*time.Second))
	//c.Get("a") // 1, true
	//c.Get("b") // nil, false
}
