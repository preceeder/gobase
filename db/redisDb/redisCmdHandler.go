/*
File Name:  redisCmdHandler.py
Description:
Author:      Chenghu
Date:       2023/8/18 17:09
Change Activity:
*/
package redisDb

import (
	"github.com/preceeder/gobase/utils"
	"strings"
)

// 这里的 any 可以是 数组(需要解开)， string, int
func RedisCmdBindName(str string, args map[string]any) ([]any, error) {
	//str := "hello, {{bushi}} huusd{} {{hs}} dd"
	//sr := map[string]any{"bushi": "23", "hs": []int{1, 2, 3}}
	var query = make([]any, 0)
	var temP = strings.Split(str, " ")
	for _, v := range temP {
		v = strings.Trim(v, " ")
		if v == "" {
			continue
		}
		start := strings.Index(v, "{{")
		end := strings.Index(v, "}}")
		if start > 0 && end < len(v) { // 名字和其他的字符串有链接的只能是一个字符串
			nameSt, err := utils.StrBindName(v, args, []byte(""))
			if err != nil {
				panic("redisDb cmd err：")
			}
			query = append(query, nameSt)
		} else if start == -1 || end == -1 {
			query = append(query, v)
			continue
		} else { // 单个的 名字{{value}}
			name := v[start+2 : end]
			nameSt, err := utils.AnyToSlice(args[name]) // "name:{{nihas}}"
			if err != nil {
				panic("redisDb cmd err：")
			}
			query = append(query, nameSt...)
		}
	}
	return query, nil
}
