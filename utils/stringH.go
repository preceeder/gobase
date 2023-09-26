package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// 随机字符串
var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var lettersInt = []rune("0123456789")

func RandStr(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(letters))]
	}
	return string(rand_bytes)
}

func RandStrInt(str_len int) string {
	rand_bytes := make([]rune, str_len)
	for i := range rand_bytes {
		rand_bytes[i] = letters[rand.Intn(len(lettersInt))]
	}
	return string(rand_bytes)
}

// 最少13位
func GenterWithoutRepetitionStr(strl int) string {
	tim := strconv.FormatInt(time.Now().UnixMilli(), 10)
	if strl < 13 {
		return tim
	}
	netstr := RandStr(strl - 13)
	return tim + netstr
}

func GenterWithoutRepetitionInt(strl int) string {
	tim := strconv.FormatInt(time.Now().UnixMilli(), 10)
	if strl < 13 {
		return tim
	}
	netstr := RandStrInt(strl - 13)
	return tim + netstr
}

// 这里的 any 可以是 数组(需要解开)， string, int

func StrBindName(str string, args map[string]any, spacing []byte) (temPs string, err error) {
	//str := "hello, {{bushi}} huusd{} {{hs}} dd"
	//sr := map[string]any{"bushi": "23", "hs": []int{1, 2, 3}}
	temPs = ""
	for {
		start := strings.Index(str, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(str, "}}")
		if end == -1 {
			break
		}
		name := str[start+2 : end]
		nameSt, _ := AnyToString(args[name], spacing)
		temPs += str[:start] + nameSt
		str = str[end+2:]
	}
	return temPs + str, nil
	//fmt.Println(temPs)
}
