/*
File Name:  goRedis_test.go.go
Description:
Author:      Chenghu
Date:       2023/8/30 12:00
Change Activity:
*/
package redisDb

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

var redisObj = map[string]RedisConfig{
	"default": RedisConfig{
		Host:        "127.0.0.1",
		Port:        "5637",
		Password:    "QDjk9UdkoD6cv",
		MaxIdle:     2,
		IdleTimeout: 240,
		PoolSize:    10,
		Db:          0,
	},
}

func TestDo(t *testing.T) {
	type args struct {
		cmd  map[string]any
		agrs map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Do(tt.args.cmd, tt.args.agrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Do() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecScript(t *testing.T) {
	type args struct {
		sc     map[string]any
		keys   map[string]string
		values map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "", args: args{
			sc: map[string]any{
				"db":      "default",
				"script":  "",
				"keys":    []string{"room_id"},
				"argv":    []string{"which", "who", "room_id_timestamp"},
				"default": map[string]any{"which": 123, "who": 234},
			},
			keys:   map[string]string{"room_id": "sdwww"},
			values: map[string]any{"room_id_timestamp": "sdff"},
		}},
	}
	InitRedisWithStruct(redisObj)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecScript(tt.args.sc, tt.args.keys, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecScript() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecScript() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitRedisWithStruct(t *testing.T) {
	type args struct {
		config map[string]RedisConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRedisWithStruct(tt.args.config)
		})
	}
}

func TestInitRedisWithViperConfig(t *testing.T) {
	type args struct {
		config viper.Viper
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitRedisWithViperConfig(tt.args.config)
		})
	}
}

func TestRedisClose(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RedisClose()
		})
	}
}

func TestScript(t *testing.T) {
	type args struct {
		db     *redis.Client
		script string
		keys   []string
		values []any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Script(tt.args.db, tt.args.script, tt.args.keys, tt.args.values)
			if (err != nil) != tt.wantErr {
				t.Errorf("Script() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Script() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDb(t *testing.T) {
	type args struct {
		d *string
	}
	tests := []struct {
		name    string
		args    args
		wantDb  *redis.Client
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDb, err := getDb(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDb() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDb, tt.wantDb) {
				t.Errorf("getDb() gotDb = %v, want %v", gotDb, tt.wantDb)
			}
		})
	}
}

func Test_getScriptKv(t *testing.T) {
	type args[T interface{ string | any }] struct {
		mm     map[string]any
		kmm    string
		source map[string]T
		fData  map[string]any
	}
	type testCase[T interface{ string | any }] struct {
		name    string
		args    args[T]
		want    []T
		wantErr bool
	}
	tests := []testCase[string]{
		// TODO: Add test cases.
		{name: "", args: args[string]{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getScriptKv(tt.args.mm, tt.args.kmm, tt.args.source, tt.args.fData)
			if (err != nil) != tt.wantErr {
				t.Errorf("getScriptKv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScriptKv() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initRedis(t *testing.T) {
	type args struct {
		config map[string]RedisConfig
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initRedis(tt.args.config)
		})
	}
}

func Test_readRedisConfig(t *testing.T) {
	type args struct {
		v viper.Viper
	}
	tests := []struct {
		name   string
		args   args
		wantRs map[string]RedisConfig
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRs := readRedisConfig(tt.args.v); !reflect.DeepEqual(gotRs, tt.wantRs) {
				t.Errorf("readRedisConfig() = %v, want %v", gotRs, tt.wantRs)
			}
		})
	}
}
