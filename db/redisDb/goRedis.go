/*
File Name:  goRedis.py
Description:
Author:      Chenghu
Date:       2023/8/18 13:39
Change Activity:

	var redisObj = map[string]RedisConfig{
		"default": RedisConfig{
			Host:        "127.0.0.1",
			Port:        "5637",
			Password:    "",
			MaxIdle:     2,
			IdleTimeout: 240,
			PoolSize:    10,
			Db:          0,
		},
	}
*/
package redisDb

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/preceeder/gobase/utils"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"reflect"
	"strconv"
	"time"
)

var Rd map[string]*redis.Client

type RedisConfig struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	Db          int    `json:"db"`
	MaxIdle     int    `json:"maxIdle"`
	IdleTimeout int    `json:"idleTimeout"`
	PoolSize    int    `json:"PoolSize"`
}

type ExpTime func() int64

/*
var redisObj = map[string]RedisConfig{
	"default": RedisConfig{
		Host:        "127.0.0.1",
		Port:        "5637",
		Password:    "",
		MaxIdle:     2,
		IdleTimeout: 240,
		PoolSize:    10,
		Db:          0,
	},
}

*/

//func init(){
// // configObj 的 viper 变量没有弄出来
//	config := *gobase.ConfigObj.viper
//	redisConfig := readRedisConfig(config)
//	initRedis(redisConfig)
//}

// 使用 viper读取的配置初始化
func InitRedisWithViperConfig(config viper.Viper) {
	redisConfig := readRedisConfig(config)
	initRedis(redisConfig)
}

func readRedisConfig(v viper.Viper) (rs map[string]RedisConfig) {
	redisDb := v.Sub("redis")
	if redisDb == nil {
		fmt.Printf("redisDb config is nil")
		os.Exit(1)
	}
	rs = make(map[string]RedisConfig)
	err := redisDb.Unmarshal(&rs)
	if err != nil {
		fmt.Printf("redisDb config read error: " + err.Error())
		os.Exit(1)
	}
	return
}

// 使用传入的 redisDb 配置文件初始化

func InitRedisWithStruct(config map[string]RedisConfig) {
	initRedis(config)
}

func initRedis(config map[string]RedisConfig) {
	Rd = make(map[string]*redis.Client)
	for key, v := range config {
		rcs, _ := sonic.MarshalString(v)
		slog.Info("redisDb connect: " + rcs)
		addr := v.Host + ":" + v.Port
		redisOpt := &redis.Options{
			Addr:         addr,
			Password:     v.Password,
			DB:           v.Db,
			PoolSize:     v.PoolSize,
			MaxIdleConns: v.MaxIdle,
			MinIdleConns: v.MaxIdle,
		}
		Rd[key] = redis.NewClient(redisOpt)
		Rd[key].AddHook(RKParesHook{})
	}
}

func RedisClose() {
	if Rd != nil {
		for key, v := range Rd {
			err := v.Close()
			str := ""
			if err != nil {
				str, _ = utils.StrBindName("close redisDb {{index}} {{error}}", map[string]any{"index": key, "error": err.Error()}, []byte(""))
			} else {
				str, _ = utils.StrBindName("close redisDb {{index}}", map[string]any{"index": key}, []byte(""))

			}
			slog.Info(str)
		}
	}
}

func getScriptKv[T string | any](context2 utils.Context, mm map[string]any, kmm string, source map[string]T, fData map[string]any) ([]T, error) {
	var sk = make([]T, 0)
	if kd, ok := mm[kmm]; ok {
		vv := kd.([]string)
		for _, kddd := range vv {
			if kv, ok := source[kddd]; ok {
				sk = append(sk, kv)
			} else if kv, ok := fData[kddd]; ok {
				sk = append(sk, kv.(T))
			} else {
				slog.Error("ExecScript error", "requestId", context2.RequestId, "error", "keys key is not enough", "key", kddd)
				return nil, utils.BaseRunTimeError{ErrorCode: 500, Message: "redis error"}
			}
		}
	}
	return sk, nil
}

func getDb(d *string) (db *redis.Client, err error) {
	if d != nil {
		if dd, ok := Rd[*d]; !ok {
			err = utils.BaseRunTimeError{ErrorCode: 500, Message: "redis db" + *d + " 不存在"}
		} else {
			db = dd
		}
	} else {
		db = Rd["default"]
	}
	return
}

// # 删除 通话缓存记录
//
//	DEL_CALLING_RECORD = {
//	   "db": “default”,
//	   "script": DEL_CALLING_RECORD_LUA,
//	   "keys": []string{"room_id"},
//	   "argv": []string{"which", "who", "room_id_timestamp"},
//	   "default": map[string]any
//	}

// 执行 lua 脚本
func ExecScript(context2 utils.Context, sc map[string]any, keys map[string]string, values map[string]any) (*redis.Cmd, error) {
	tdb := sc["db"].(string)
	rd, err := getDb(&tdb)
	if err != nil {
		slog.Error(err.Error(), "requestId", context2.RequestId)
		panic(err)
	}
	var fData = make(map[string]any)
	if fd, ok := sc["default"]; ok {
		fdd := fd.(map[string]any)
		for fk, fv := range fdd {
			var td any
			if fn, ok := fv.(ExpTime); ok {
				td = fn()
			} else {
				td = fv
			}
			fData[fk] = td
		}

	}

	sk, err := getScriptKv(context2, sc, "keys", keys, fData)
	if err != nil {
		panic(utils.BaseRunTimeError{ErrorCode: 500, Message: "redis error1"})
	}
	argvs, err := getScriptKv(context2, sc, "argv", values, fData)
	if err != nil {
		panic(utils.BaseRunTimeError{ErrorCode: 500, Message: "redis error2"})
	}
	res, err := Script(context2, rd, sc["script"].(string), sk, argvs)
	if err != nil {
		panic(utils.BaseRunTimeError{ErrorCode: 500, Message: "redis error3"})
	}
	return res, nil
}

func Script(context2 utils.Context, db *redis.Client, script string, keys []string, values []any) (*redis.Cmd, error) {
	//script := `
	//local key1 = KEYS[1]
	//local user_id = ARGV[1]
	//local result_data = {{}, {}}
	//result_data[1] = redisDb.call("set", key1, user_id)
	//result_data[2] = redisDb.call("GET", key1)
	//return result_data`
	ctx := context.Background()
	hesHasScript := cryptor.Sha1(script)
	var rd *redis.Client
	if db == nil {
		rd = Rd["default"]
	} else {
		rd = db
	}
	ns, err := rd.ScriptExists(ctx, hesHasScript).Result()
	if err != nil {
		slog.Error("ScriptExists error", "error", err.Error(), "requestId", context2.RequestId)
		return nil, err
	}
	//不存在则加载
	if !ns[0] {
		_, err := rd.ScriptLoad(ctx, script).Result()
		if err != nil {
			slog.Error("ScriptLoad error", "error", err.Error(), "requestId", context2.RequestId)
			return nil, err
		}
	}
	rcmd := rd.EvalSha(ctx, hesHasScript, keys, values...)
	err = rcmd.Err()
	if err != nil {
		slog.Error("EvalSha error", "error", err.Error(), "requestId", context2.RequestId)
		return nil, err
	}
	return rcmd, nil
}

func Do(context2 utils.Context, cmd map[string]any, agrs map[string]any, includeArgs ...any) (*redis.Cmd, error) {
	/*
			cmd   map[string]any = {
					"cmd": "get {{ss}}
		            "includeArgs": []any   // 附加的参数
					"key": "sdf{{}}"   // 设置 exp的时候需要用到
					"exp": 234｜ time.Duration｜func()time.Duration{return 2},      // s
					"db": "",
					}
	*/
	ctx := context.Background()
	cmdStr, err := RedisCmdBindName(cmd["cmd"].(string), agrs)
	if err != nil {
		slog.Error("redisDb args pares err", "args", agrs, "cmd", cmd["cmd"], "requestId", context2.RequestId)
		panic("redisDb args pares err")
	}

	var rd *redis.Client
	if willDb, ok := cmd["db"]; ok {
		rd = Rd[willDb.(string)]
	} else {
		rd = Rd["default"]
	}
	if len(includeArgs) > 0 {
		cmdStr = append(cmdStr, includeArgs...)
	}

	rcmd := rd.Do(ctx, cmdStr...)
	if rcmd.Err() != nil {
		slog.Error("redisDb exec failed", "error", rcmd.Err().Error(), "requestId", context2.RequestId)
		return nil, rcmd.Err()
	}

	// 设置过期时间
	if exp, exist := cmd["exp"]; exist {
		// int 秒, fun()
		var expt time.Duration = 0
		switch exp.(type) {
		case int64, int, int32:
			expIn64, _ := convertor.ToInt(exp)
			expt = time.Duration(expIn64) * time.Second
		case time.Duration:
			expt = exp.(time.Duration)
		case func() time.Duration:
			vd := reflect.ValueOf(exp)
			ex := vd.Call(nil)
			expt = ex[0].Interface().(time.Duration)
		}

		rkey, ok := cmd["key"]
		if !ok {
			slog.Error("set exp, cmd must has key", "cmd", cmd["cmd"], "requestId", context2.RequestId)
			panic("set exp, cmd must has key")
		}
		keyStr, err := utils.StrBindName(rkey.(string), agrs, []byte(""))
		if err != nil {
			slog.Error("redisDb key args pares err", "requestId", context2.RequestId)
			panic("redisDb key args pares err")
		}
		err = rd.Expire(ctx, keyStr, expt).Err()
		if err != nil {
			slog.Error("redisDb key set exp err: "+err.Error(), "requestId", context2.RequestId)
			panic("redisDb key set exp err: " + err.Error())
		}
	}

	return rcmd, nil
}

// 分布式事务锁
type RedisLocks struct {
	Value    string
	Key      string
	Exp      int64
	db       *redis.Client
	Context2 utils.Context
}

// pao一个lua 脚本  加锁， 释放锁
func (r *RedisLocks) GetLock(db ...string) bool {
	// exp  ms
	if len(db) == 1 {
		r.db = Rd[db[0]]
	} else {
		r.db = Rd["default"]
	}
	ctx := context.Background()
	if r.Value == "" {
		value := strconv.FormatInt(time.Now().UnixMilli(), 10) + "_" + utils.RandStr(5)
		r.Value = value
	}

	if r.Exp == 0 {
		r.Exp = 5
	}

	expl := time.Duration(r.Exp * 1000000)
	res, err := r.db.SetNX(ctx, r.Key, r.Value, expl).Result()
	if err != nil {
		slog.Error("redisDb 加锁失败 key: "+r.Key+"err: "+err.Error(), "requestId", r.Context2.GetString("requestId"))
		panic("redisDb 加锁失败 key: " + r.Key + "err: " + err.Error())
	}
	return res
}

func (r *RedisLocks) ReleaseLock() {
	// 删除锁
	var del = `
	if redis.call("get",KEYS[1]) == ARGV[1] then
		return redisDb.call("del",KEYS[1])
	else
		return 0
	end`
	ok, err := Script(r.Context2, r.db, del, []string{r.Key}, []any{r.Value})
	if err != nil {
		panic("redisDb ReleaseLock err: " + err.Error())
	}
	fmt.Println(ok)
}
