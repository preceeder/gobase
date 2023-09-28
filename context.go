/*
File Name:  context.py
Description:
Author:      Chenghu
Date:       2023/9/28 14:15
Change Activity:
*/
package gobase

import (
	"context"
	"time"
)

func sud(ctx context.Context) {
}

type Context struct {
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

func sde() {
	sud(Context{})
}
