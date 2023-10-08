/*
File Name:  context.py
Description:
Author:      Chenghu
Date:       2023/9/28 14:15
Change Activity:
*/
package utils

import (
	"fmt"
	"sync"
	"time"
)

type Context struct {
	m *sync.Map
}

func (y Context) Deadline() (deadline time.Time, ok bool) {
	//TODO implement me
	panic("implement me")
}

func (y Context) Done() <-chan struct{} {
	//TODO implement me
	panic("implement me")
}

func (y Context) Err() error {
	//TODO implement me
	panic("implement me")
}

func (y Context) Value(key any) any {
	//TODO implement me
	panic("implement me")
}

func (c *Context) Set(key any, v any) {
	if c.m == nil {
		c.m = new(sync.Map)
	}
	c.m.Store(key, v)
}

func (c *Context) GetString(key any) (value string) {
	val, ok := c.Get(key)
	if ok && val != nil {
		value, _ = val.(string)
	}
	return
}

func (c *Context) Get(key any) (value any, exists bool) {
	value, exists = c.m.Load(key)
	return
}

func sde() {
	con := Context{}
	con.Set("requestId", "ok")
	fmt.Printf(con.GetString("requestId"))
}
