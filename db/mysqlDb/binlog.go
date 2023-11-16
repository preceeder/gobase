/*
File Name:  binlog.go
Description:
Author:      Chenghu
Date:       2023/10/13 15:23
Change Activity:
*/

// use

// 1. 调用 SetBinlogTable(),  配置表， 以及表数据怎样处理

// 2 有配置文件   配置文件的写法
//
//	"binlog": {
//		 "addr": "host:port",
//		 "password": "xxxxxx",
//		 "user": "xxx",
//		 "db": "xxx"
//	}
// 然后直接 调用 Run()

// 2 没有配置文件 就在 调用 Run(binlogConfig) 函数时 传入实例化的 BinlogConfig

package mysqlDb

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/preceeder/gobase/utils"
	"github.com/spf13/viper"
	"log/slog"
	"os"
	"reflect"
)

type EventHandler struct {
	canal.DummyEventHandler
}

type BinlogConfig struct {
	Addr     string `json:"addr"` // "127.0.0.1:13306"
	Password string `json:"password"`
	User     string `json:"user"`
	Db       string `json:"db"`
}

var BConfig BinlogConfig

// 使用 viper读取的配置初始化
func InitBinlogWithViperConfig(config viper.Viper) {
	//BConfig = readBinlogConfig(config)
	utils.ReadViperConfig(config, "binlog", &BConfig)
}

//func readBinlogConfig(config viper.Viper) BinlogConfig {
//	binlog := config.Sub("binlog")
//	if binlog == nil {
//		fmt.Printf("binlog config is nil")
//		os.Exit(1)
//	}
//	var subbinlog BinlogConfig
//	err := binlog.Unmarshal(&subbinlog)
//	if err != nil {
//		fmt.Printf("subbinlog config read error: " + err.Error())
//		os.Exit(1)
//	}
//	return subbinlog
//}

type TableData struct {
	TableName  string
	Table      reflect.Type
	Action     []string // 需要监听的事件 update, insert, delete
	HandlerFun func(action string, oldData any, newData ...any)
}

var tableData = make(map[string]TableData, 0)

// tableName 表名，  data 需要的 表结构体
// HandlerFun
//  action=insert 只有一个参数
//  action=update 有两个参数
//  action=delete 只有一个参数

func SetBinlogTable(tableName string, action []string, data any, HandlerFun func(action string, oldData any, newData ...any)) {
	rd := reflect.Indirect(reflect.ValueOf(data))
	if rd.Type().Kind() != reflect.Struct {
		panic("table type error")
	}
	tableData[tableName] = TableData{TableName: tableName, Action: action, Table: reflect.TypeOf(data), HandlerFun: HandlerFun}
}

func (h *EventHandler) OnRow(e *canal.RowsEvent) error {
	v, ok := tableData[e.Table.Name]
	if !ok {
		return nil
	}
	if !slice.Contain(v.Action, e.Action) {
		return nil
	}
	var EnumMap map[int][]string = make(map[int][]string)

	columsName := make([]string, len(e.Table.Columns))
	for i, v := range e.Table.Columns {
		columsName[i] = v.Name
		if v.Type == 3 {
			EnumMap[i] = make([]string, 0)
			EnumMap[i] = append(append(EnumMap[i], ""), v.EnumValues...)
		}
	}
	targetData := make([]any, len(e.Rows))
	for i, edata := range e.Rows {
		for k, v := range EnumMap {
			if ev := edata[k]; ev != nil {
				edata[k] = v[ev.(int64)]
			}
		}
		data, _ := utils.SliceToMap(columsName, edata)
		//dd, _ := sonic.Marshal(data)   // 使用json 处理 []uint8 的数据会有问题
		if ok {
			table := reflect.New(v.Table).Interface()
			err := utils.MapToStructWithTag(data, table, "json", true)
			//err := sonic.UnmarshalString(dd, table)
			if err != nil {
				slog.Error("binlog error", "error", err.Error())
			}
			targetData[i] = table
		}
	}
	// 保存数据 操作 由用户自己决定
	go v.HandlerFun(e.Action, targetData[0], targetData[1:]...)

	return nil
}

func (h *EventHandler) String() string {
	return "MyEventHandler"
}

// 在调用这个函数之前 必须先调用  SetBinlogTable  设置需要监听的 表， 一旦启动就没发在设置了

func BinLogRun(config ...BinlogConfig) {
	var binconfig BinlogConfig
	if len(config) == 0 {
		binconfig = BConfig
	} else {
		binconfig = config[0]
	}
	cfg := canal.NewDefaultConfig()
	cfg.Addr = binconfig.Addr
	cfg.Password = binconfig.Password
	cfg.User = binconfig.User
	cfg.Charset = "utf8mb4"
	cfg.Dump.ExecutionPath = ""
	// We only care table canal_test in test db
	cfg.Dump.TableDB = binconfig.Db
	cfg.Dump.Databases = []string{binconfig.Db}
	for _, v := range tableData {
		cfg.Dump.Tables = append(cfg.Dump.Tables, v.TableName)
		cfg.IncludeTableRegex = append(cfg.IncludeTableRegex, fmt.Sprintf("^%s\\.%s$", binconfig.Db, v.TableName))
	}
	slog.Info("binglog config", "config", cfg)
	c, err := canal.NewCanal(cfg)

	if err != nil {
		slog.Error(err.Error())
	}

	c.SetEventHandler(&EventHandler{})

	//开启信号监听
	signl := utils.StartSignalLister()

	//开启信号处理
	go utils.SignalHandler(signl, func() {
		//平滑关闭
		os.Exit(1)
	})
	slog.Info("开启 binglo 监听", "config", binconfig)
	// Start canal
	pos, err := c.GetMasterPos()
	if err != nil {
		slog.Error("获取binlog 最新日志位置失败")
	}
	c.RunFrom(pos)
}
